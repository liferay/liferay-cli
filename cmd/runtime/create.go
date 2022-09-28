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
			Verbose, &wg)

		wg.Add(1)
		go lcectldocker.BuildImage("localdev-server", path.Join(
			viper.GetString(constants.Const.RepoDir), "docker", "images", "localdev-server"),
			Verbose, &wg)

		wg.Wait()

		s.Stop()
		s.Suffix = " Creating localdev environment..."
		s.FinalMSG = fmt.Sprintf("\u2705 Created localdev environment.\n")
		s.Restart()

		wg.Add(1)
		lcectldocker.InvokeCommandInLocaldev(
			"localdev-start", []string{"/repo/scripts/cluster-start.sh"}, Verbose, &wg)

		wg.Wait()
		s.Stop()

		return nil
	},
}

func init() {
	createCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose output")
	runtimeCmd.AddCommand(createCmd)
}
