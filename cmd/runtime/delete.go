/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package runtime

import (
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"liferay.com/lcectl/constants"
	lcectldocker "liferay.com/lcectl/docker"
	"liferay.com/lcectl/git"
)

// createCmd represents the create command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the runtime environment for Liferay Client Extension development",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Color("green")
		s.Suffix = " Synchronizing localdev sources..."
		s.FinalMSG = fmt.Sprintf("\u2705 Synced localdev sources.\n")
		s.Start()

		git.SyncGit()

		s.Stop()
		s.Suffix = " Building localdev image..."
		s.FinalMSG = fmt.Sprintf("\u2705 Built localdev images.\n")
		s.Restart()

		var wg sync.WaitGroup
		wg.Add(1)
		go lcectldocker.BuildImage("localdev-server", path.Join(
			viper.GetString(constants.Const.RepoDir), "docker", "images", "localdev-server"),
			Verbose, &wg)

		wg.Wait()

		s.Stop()
		s.Suffix = " Deleting localdev environment..."
		s.FinalMSG = fmt.Sprintf("\u2705 Deleted localdev environment.\n")
		s.Restart()

		wg.Add(1)

		config := container.Config{
			Image: "localdev-server",
			Cmd:   []string{"/repo/scripts/cluster-delete.sh"},
		}
		host := container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
				"/var/run/docker.sock:/var/run/docker.sock",
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		}

		lcectldocker.InvokeCommandInLocaldev("localdev-delete", config, host, Verbose, &wg)

		wg.Wait()
		s.Stop()

		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose output")
	runtimeCmd.AddCommand(deleteCmd)
}
