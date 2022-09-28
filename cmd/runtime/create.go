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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"liferay.com/lcectl/constants"
	lcectldocker "liferay.com/lcectl/docker"
	"liferay.com/lcectl/git"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates the runtime environment for Liferay Client Extension development",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dockerClient, err := lcectldocker.GetDockerClient()

		if err != nil {
			return err
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Color("green")
		s.Suffix = " Synchronizing localdev sources..."
		s.FinalMSG = fmt.Sprintf("\u2705 Synced localdev sources.\n")
		s.Start()

		git.SyncGit()

		s.Stop()
		s.Suffix = " Building localdev images..."
		s.FinalMSG = fmt.Sprintf("\u2705 Built localdev images.\n")
		s.Restart()

		var wg sync.WaitGroup
		wg.Add(1)
		go lcectldocker.BuildImage("dxp-server", path.Join(
			viper.GetString(constants.Const.RepoDir), "docker", "images", "dxp-server"),
			dockerClient, &wg)

		wg.Add(1)
		go lcectldocker.BuildImage("localdev-server", path.Join(
			viper.GetString(constants.Const.RepoDir), "docker", "images", "localdev-server"),
			dockerClient, &wg)

		wg.Wait()

		s.Stop()
		s.Suffix = " Creating localdev environment..."
		s.FinalMSG = fmt.Sprintf("\u2705 Created localdev environment.\n")
		s.Restart()

		wg.Add(1)
		lcectldocker.InvokeCommandInLocaldev("localdev-start", []string{"/repo/scripts/cluster-start.sh"}, dockerClient, &wg)

		wg.Wait()
		s.Stop()

		return nil
	},
}

func init() {
	runtimeCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
