/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"liferay.com/liferay/cli/ansicolor"
	"liferay.com/liferay/cli/constants"
	"liferay.com/liferay/cli/docker"
	"liferay.com/liferay/cli/ext"
	"liferay.com/liferay/cli/flags"
	"liferay.com/liferay/cli/spinner"
)

var argRegex = regexp.MustCompile("^--(.*)=(.*)$")
var noPrompt bool
var whitespace = regexp.MustCompile(`\s`)

type void struct{}

var novalue void

var createCmd = &cobra.Command{
	Use:   "create [OPTIONS] [FLAGS]",
	Short: "Creates new Client Extensions using a wizard-like interface",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if noPrompt {
			invokeCreate(cmd, args)
		}

		actionIdx, _ := selection("How you would like to proceed", []string{
			ansicolor.Bold("Create") + " project from " + ansicolor.Bold("sample"),
			ansicolor.Bold("Create") + " project from " + ansicolor.Bold("template"),
			ansicolor.Bold("Modify") + " existing project",
		})

		switch actionIdx {
		case 0:
			createFrom(cmd, "sample")
		case 1:
			createFrom(cmd, "template")
		case 2:
			modifyExistingProject(cmd)
		}
	},
}

func applyPartialTo(cmd *cobra.Command, project string) {
	resourceType := "partial"
	resources := getTypeSubset(resourceType, getClientExtentionResourcesJson())

	// TODO filter by matching type

	// Get the template type from the project
	yaml := readClientExtensionYamlFromProject(project)
	runtime := yaml["runtime"].(map[interface{}]interface{})
	projectTemplate := runtime["template"].(string)

	partials := make(map[string]map[string]interface{})
	for key, partial := range resources {
		if partial["template"] == projectTemplate {
			partials[key] = partial
		}
	}

	var resource map[string]interface{}

	if len(partials) > 0 {
		resource = selectResourceByName(resourceType, partials)

		createFromResource(cmd, resource, project)
	} else {
		fmt.Println("There are no partials that apply to", project)
	}
}

func assembleSelectKey(template map[string]interface{}) string {
	key := ansicolor.Bold(template["name"].(string))
	if template["description"] != nil {
		key += " - " + template["description"].(string)
	}
	return key
}

func createFrom(cmd *cobra.Command, resourceType string) {
	resources := getTypeSubset(resourceType, getClientExtentionResourcesJson())
	categories := getCategoriesFromSubset(resources)
	var resource map[string]interface{}

	if len(categories) > 1 {
		listByIdx, _ := selection(fmt.Sprintf("List %ss by", resourceType), []string{
			ansicolor.Bold("Category"),
			ansicolor.Bold("Name"),
		})

		switch listByIdx {
		case 0:
			resource = selectResourceByName(resourceType, listByCategory(resourceType, categories))
		case 1:
			resource = selectResourceByName(resourceType, resources)
		}
	} else {
		resource = selectResourceByName(resourceType, resources)
	}

	createFromResource(cmd, resource, promptForWorkspacePath(resource["name"].(string)))
}

func createFromResource(cmd *cobra.Command, resource map[string]interface{}, workspacePath string) {
	args := make([]interface{}, 0)
	if resource["args"] != nil {
		args = resource["args"].([]interface{})
	}

	generatorArgs := make([]string, len(args)+2)
	generatorArgs[0] = fmt.Sprintf("--resource-path=%s", resource["type"].(string)+"/"+resource["name"].(string))
	generatorArgs[1] = fmt.Sprintf("--workspace-path=%s", workspacePath)
	var argIdx = 2

	for _, arg := range args {
		argEntry := (arg).(map[string]interface{})
		argDefault := ""

		if argEntry["default"] != nil {
			argDefault = argEntry["default"].(string)
		}

		argName := argEntry["name"].(string)

		value := prompt(
			fmt.Sprintf(argEntry["description"].(string)),
			fmt.Sprintf("Specify '%s'", argName),
			argDefault,
			func(input string) error {
				if len(input) <= 0 {
					return errors.New(argName + " must not be empty")
				}
				return nil
			},
		)

		generatorArgs[argIdx] = fmt.Sprintf("--args=%s=%s", argEntry["name"].(string), value)
		argIdx++
	}

	idx, _ := selection(
		"Ready to finish",
		[]string{
			"Create",
			"Just show the command",
		},
	)

	switch idx {
	case 0:
		invokeCreate(cmd, generatorArgs)
	case 1:
		cmd := "liferay ext create --noprompt --"
		for _, garg := range generatorArgs {
			cmd += argRegex.ReplaceAllString(garg, " --$1=\"$2\"")
		}
		fmt.Println(cmd)
	}
}

func getCategoriesFromSubset(subset map[string]map[string]interface{}) map[string]map[string]map[string]interface{} {
	categories := make(map[string]map[string]map[string]interface{})

	for _, entry := range subset {
		category := entry["category"]

		if category == nil {
			category = "General"
		}

		if categories[category.(string)] == nil {
			categories[category.(string)] = make(map[string]map[string]interface{}, 0)
		}

		categories[category.(string)][assembleSelectKey(entry)] = entry
	}

	return categories
}

func getClientExtentionResourcesJson() []map[string]interface{} {
	bytes, err := os.ReadFile(filepath.Join(viper.GetString(constants.Const.RepoDir), "resources", "client-extension-resources.json"))

	if err != nil {
		panic(err)
	}

	var data []map[string]interface{}

	if err := json.Unmarshal(bytes, &data); err != nil {
		panic(err)
	}

	return data
}

func getWorkspaceProjects() []string {
	workspaceDir := viper.GetString(constants.Const.ExtClientExtensionDir)
	projectSet := make(map[string]void)

	e := filepath.Walk(
		workspaceDir,
		func(path string, info os.FileInfo, err error) error {
			if err == nil && info.Name() == "client-extension.yaml" {
				fullPathDir := filepath.Dir(path)
				projectSet[fullPathDir[len(workspaceDir):]] = novalue
			}
			return nil
		},
	)

	if e != nil {
		log.Fatal(e)
	}

	projects := make([]string, 0)
	for key := range projectSet {
		projects = append(projects, key)
	}

	return projects
}

func getTypeSubset(subsetType string, clientExtentionResources []map[string]interface{}) map[string]map[string]interface{} {
	count := 0

	for _, fi := range clientExtentionResources {
		if fi["type"] != nil && fi["type"].(string) != subsetType {
			continue
		}
		count++
	}

	subset := make(map[string]map[string]interface{}, count)

	for _, entry := range clientExtentionResources {
		if entry["type"] != nil && entry["type"].(string) != subsetType {
			continue
		}
		subset[assembleSelectKey(entry)] = entry
	}

	return subset
}

func init() {
	extCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVarP(&noPrompt, "noprompt", "n", false, "Do not show the wizard prompts, just use args")
}

func invokeCreate(cmd *cobra.Command, args []string) {
	config := container.Config{
		Image: "localdev-server",
		Cmd:   []string{"/repo/scripts/ext/create.py"},
		Env: []string{
			"CLIENT_EXTENSION_DIR_KEY=" + ext.GetExtensionDirKey(),
			"WORKSPACE_BASE_PATH=/workspace/client-extensions",
			"LOCALDEV_REPO=/repo",
			"CREATE_ARGS=" + strings.Join(args, "|"),
		},
	}
	host := container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
			docker.GetDockerSocket() + ":/var/run/docker.sock",
			fmt.Sprintf("%s:/workspace/client-extensions", flags.ClientExtensionDir),
		},
		NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
	}
	docker.PerformOSSpecificAdjustments(&config, &host)

	exitCode := spinner.Spin(
		spinner.SpinOptions{
			Doing: "Creating", Done: "created", On: "'localdev' client extension project", Enable: !flags.Verbose,
		},
		func(fior func(io.ReadCloser, bool, string) int) int {
			return docker.InvokeCommandInLocaldev("localdev-ext-create", config, host, true, flags.Verbose, fior, "")
		})

	if exitCode == 0 && runtime.GOOS == "windows" {
		refreshCmd.Run(cmd, args)
	}

	os.Exit(exitCode)
}

func listByCategory(resourceType string, categories map[string]map[string]map[string]interface{}) map[string]map[string]interface{} {
	keys := make([]string, 0)
	for key := range categories {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	_, categoryKey := selection("Choose a category", keys)

	return categories[categoryKey]
}

func modifyExistingProject(cmd *cobra.Command) {
	projects := getWorkspaceProjects()

	if len(projects) < 1 {
		fmt.Println("There are no projects.")
		os.Exit(0)
	}

	_, project := selection("Select the project", projects)

	actionIdx, _ := selection("What kind of modification will you perform", []string{
		"Apply a " + ansicolor.Bold("partial"),
		"Add a " + ansicolor.Bold("client extension"),
	})

	switch actionIdx {
	case 0:
		applyPartialTo(cmd, project)
	case 1:
		//addClientExtensionTo(project)
	}
}

func prompt(question string, label string, dflt string, validate func(input string) error) string {
	fmt.Println(question)
	prompt := promptui.Prompt{
		Label:    label,
		Default:  dflt,
		Validate: validate,
	}

	answer, err := prompt.Run()

	if err != nil {
		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}
		log.Fatal(err)
	}

	return answer
}

func promptForWorkspacePath(dflt string) string {
	return prompt(
		"Specify the workspace directory where the project should be created",
		"Directory",
		dflt,
		func(input string) error {
			if len(input) <= 0 {
				return errors.New("the directory name must not be empty")
			}
			if whitespace.MatchString(input) {
				return errors.New("the directory name must not contain spaces")
			}
			path := filepath.Join(viper.GetString(constants.Const.ExtClientExtensionDir), input)
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				return errors.New("the directory name already exists")
			}
			clientExtensionYaml := filepath.Join(filepath.Dir(path), "client-extension.yaml")
			if _, err := os.Stat(clientExtensionYaml); !os.IsNotExist(err) {
				return errors.New("a project cannot be created inside another project")
			}

			return nil
		})
}

func readClientExtensionYamlFromProject(project string) map[string]interface{} {
	yamlFile, err := ioutil.ReadFile(
		filepath.Join(
			viper.GetString(constants.Const.ExtClientExtensionDir),
			project,
			"client-extension.yaml",
		),
	)

	if err != nil {
		log.Fatal("Could not read project client-extension.yaml ", err)
	}

	var data map[string]interface{}

	if err = yaml.Unmarshal(yamlFile, &data); err != nil {
		log.Fatal("Could not read project client-extension.yaml ", err)
	}

	return data
}

func selection(label string, items interface{}) (int, string) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  10,
	}

	idx, answer, err := prompt.Run()

	if err != nil {
		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}
		log.Fatal(err)
	}

	return idx, answer
}

func selectResourceByName(resourceType string, resources map[string]map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0)
	for key := range resources {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	_, resourceKey := selection(fmt.Sprintf("Choose a %s", resourceType), keys)
	return resources[resourceKey]
}
