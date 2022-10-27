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
	"os/user"
	"path/filepath"
	"regexp"
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
			createFromSample()
		case 1:
			createFromTemplate()
		case 2:
			modifyExistingProject()
		}
	},
}

func clientExtentionResourcesJson() []map[string]interface{} {
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

func init() {
	extCmd.AddCommand(createCmd)
}

func invokeCreate(args []string) {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	config := container.Config{
		Image: "localdev-server",
		Cmd:   []string{"/repo/scripts/ext/create.py"},
		Env: []string{
			"WORKSPACE_BASE_PATH=/workspace/client-extensions",
			"LOCALDEV_REPO=/repo",
			"CREATE_ARGS=" + strings.Join(args, "|"),
		},
		User: user.Uid + ":" + user.Gid,
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

func createFromTemplateByName() {
	data := clientExtentionResourcesJson()

	i := 0
	for _, fi := range data {
		if fi["type"] != nil && fi["type"].(string) == "template" {
			i++
		}
	}

	keys := make([][]string, i)

	i = 0
	for _, fi := range data {
		if fi["type"] != nil && fi["type"].(string) != "template" {
			continue
		}

		keys[i] = make([]string, 2)
		keys[i][0] = fi["name"].(string)
		keys[i][1] = ""
		if fi["description"] != nil {
			keys[i][1] = fi["description"].(string)
		}
		i++
	}

	templateIdx, _ := selection("Choose a template", keys)
	workspacePath := promptForWorkspacePath(data[templateIdx]["name"].(string))

	args := data[templateIdx]["args"].([]interface{})

	generatorArgs := make([]string, len(args)+2)
	generatorArgs[0] = fmt.Sprintf("--resource-path=%s", data[templateIdx]["type"].(string)+"/"+data[templateIdx]["name"].(string))
	generatorArgs[1] = fmt.Sprintf("--workspace-path=%s", workspacePath)

	i = 2
	for _, arg := range args {
		argEntry := (arg).(map[string]interface{})
		argDefault := ""

		if argEntry["default"] != nil {
			argDefault = argEntry["default"].(string)
		}

		value := prompt(
			fmt.Sprintf(argEntry["description"].(string)),
			fmt.Sprintf("Specify '%s'", argEntry["name"].(string)),
			argDefault,
			func(input string) error {
				if len(input) <= 0 {
					return errors.New(argEntry["name"].(string) + " must not be empty")
				}
				return nil
			})
		generatorArgs[i] = fmt.Sprintf("--args=%s=%s", argEntry["name"].(string), value)
		i++
	}

	idx, _ := selection("Ready to finish", []string{"Create", "Just show the command"})

	switch idx {
	case 0:
		invokeCreate(generatorArgs)
	case 1:
		cmd := "liferay ext create "
		for _, garg := range generatorArgs {
			cmd += garg + " "
		}
		fmt.Println(cmd)
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

func createFromSample() {
	listByIdx, _ := selection("List samples by", []string{
		ansicolor.Bold("Category"),
		ansicolor.Bold("Name"),
	})

	switch listByIdx {
	case 0:
		fmt.Println("There are currently no categories. Check back soon!")
	case 1:
		fmt.Println("There are currently no samples. Check back soon!")
	}
}

func createFromTemplate() {
	listByIdx, _ := selection("List templates by", []string{
		ansicolor.Bold("Category"),
		ansicolor.Bold("Name"),
	})

	switch listByIdx {
	case 0:
		fmt.Println("There are currently no categories. Check back soon!")
	case 1:
		createFromTemplateByName()
	}
}

func modifyExistingProject() {
	fmt.Println("There are currently no projects. Check back soon!")
}
