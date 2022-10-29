/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	"liferay.com/liferay/cli/ansicolor"
	"liferay.com/liferay/cli/constants"
	"liferay.com/liferay/cli/docker"
	"liferay.com/liferay/cli/flags"
	"liferay.com/liferay/cli/spinner"
	"liferay.com/liferay/cli/user"
)

var whitespace = regexp.MustCompile(`\s`)

var createCmd = &cobra.Command{
	Use:   "create [OPTIONS] [FLAGS]",
	Short: "Creates new Client Extensions using a wizard-like interface",
	Run: func(cmd *cobra.Command, args []string) {
		actionIdx, _ := selection("How you would like to proceed", []string{
			ansicolor.Bold("Create") + " project from " + ansicolor.Bold("sample"),
			ansicolor.Bold("Create") + " project from " + ansicolor.Bold("template"),
			ansicolor.Bold("Modify") + " existing project",
		})

		switch actionIdx {
		case 0:
			createFrom("sample")
		case 1:
			createFrom("template")
		case 2:
			modifyExistingProject()
		}
	},
}

func assembleSelectKey(template map[string]interface{}) string {
	key := ansicolor.Bold(template["name"].(string))
	if template["description"] != nil {
		key += " - " + template["description"].(string)
	}
	return key
}

func createFrom(resourceType string) {
	resources := getTypeSubset(resourceType, getClientExtentionResourcesJson())
	categories := getCategoriesFromSubset(resources)

	if len(categories) > 1 {
		listByIdx, _ := selection(fmt.Sprintf("List %ss by", resourceType), []string{
			ansicolor.Bold("Category"),
			ansicolor.Bold("Name"),
		})

		switch listByIdx {
		case 0:
			listByCategory(resourceType, categories)
		case 1:
			createFromResourceByName(resourceType, resources)
		}
	} else {
		createFromResourceByName(resourceType, resources)
	}
}

func createFromResourceByName(resourceType string, resources map[string]map[string]interface{}) {
	keys := make([]string, 0)
	for key := range resources {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	_, resourceKey := selection(fmt.Sprintf("Choose a %s", resourceType), keys)
	resource := resources[resourceKey]
	workspacePath := promptForWorkspacePath(resource["name"].(string))

	args := make([]interface{}, 0)

	if resource["args"] != nil {
		args = resource["args"].([]interface{})
	}

	generatorArgs := make([]string, len(args)+2)
	generatorArgs[0] = fmt.Sprintf("--resource-path=%s", resource["type"].(string)+"/"+resource["name"].(string))
	generatorArgs[1] = fmt.Sprintf("--workspace-path=%s", workspacePath)

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

		generatorArgs = append(generatorArgs, fmt.Sprintf("--args=%s=%s", argEntry["name"].(string), value))
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
		invokeCreate(generatorArgs)
	case 1:
		cmd := "liferay ext create"
		for _, garg := range generatorArgs {
			cmd += " " + garg
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
}

func invokeCreate(args []string) {
	config := container.Config{
		Image: "localdev-server",
		Cmd:   []string{"/repo/scripts/ext/create.py"},
		Env: []string{
			"WORKSPACE_BASE_PATH=/workspace/client-extensions",
			"LOCALDEV_REPO=/repo",
			"CREATE_ARGS=" + strings.Join(args, "|"),
		},
	}
	if runtime.GOOS == "linux" {
		config.User = user.UserUidAndGuidString()
	}
	host := container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
			docker.GetDockerSocket() + ":/var/run/docker.sock",
			fmt.Sprintf("%s:/workspace/client-extensions", flags.ClientExtensionDir),
		},
		NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
	}

	exitCode := spinner.Spin(
		spinner.SpinOptions{
			Doing: "Creating", Done: "created", On: "'localdev' client extension project", Enable: !flags.Verbose,
		},
		func(fior func(io.ReadCloser, bool, string) int) int {
			return docker.InvokeCommandInLocaldev("localdev-ext-create", config, host, true, flags.Verbose, fior, "")
		})

	os.Exit(exitCode)
}

func listByCategory(resourceType string, categories map[string]map[string]map[string]interface{}) {
	keys := make([]string, 0)
	for key := range categories {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	_, categoryKey := selection("Choose a category", keys)

	createFromResourceByName(resourceType, categories[categoryKey])
}

func modifyExistingProject() {
	fmt.Println("There are currently no projects. Check back soon!")
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

func selection(label string, items interface{}) (int, string) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
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
