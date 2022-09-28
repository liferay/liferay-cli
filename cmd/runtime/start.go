/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package runtime

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"liferay.com/lcectl/constants"
	lcectldocker "liferay.com/lcectl/docker"
	"liferay.com/lcectl/prereq"
	lcectlspinner "liferay.com/lcectl/spinner"
)

// createCmd represents the create command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a stopped runtime environment for Liferay Client Extension development",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		prereq.Prereq(Verbose)

		config := container.Config{
			Image: "localdev-server",
			Cmd:   []string{"/repo/scripts/cluster-start.sh"},
		}
		host := container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
				"/var/run/docker.sock:/var/run/docker.sock",
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		}

		var s *spinner.Spinner

		if !Verbose {
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Color("green")
			s.Suffix = " Starting 'localdev' environment..."
			s.FinalMSG = fmt.Sprintf("\u2705 Started 'localdev' environment.\n")
			s.Start()
		}

		pipeSpinner := lcectlspinner.SpinnerPipe(s, " Starting 'localdev' environment [%s]", Verbose)

		signal := lcectldocker.InvokeCommandInLocaldev("localdev-start", config, host, Verbose, pipeSpinner)

		if s != nil {
			if signal > 0 {
				s.FinalMSG = fmt.Sprintf("\u2718 Something went wrong...\n")
			}

			s.Stop()
		}
	},
}

func init() {
	startCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose output")
	runtimeCmd.AddCommand(startCmd)
}
