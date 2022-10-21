/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package runtime

import (
	"github.com/spf13/cobra"
	"liferay.com/liferay/cli/flags"
	"liferay.com/liferay/cli/mkcert"
)

var install bool
var uninstall bool

// mkcertCmd represents the mkcert command
var mkcertCmd = &cobra.Command{
	Use:   "mkcert",
	Short: "Uses mkcert package to make locally-trusted development certificates.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		mkcert.VerifyRootCALoaded(flags.Verbose)

		if install {
			mkcert.InstallRootCA()
			return
		} else if uninstall {
			mkcert.UninstallRootCA()
			return
		}

		mkcert.MakeCert()
	},
}

func init() {
	runtimeCmd.AddCommand(mkcertCmd)
	mkcertCmd.Flags().BoolVarP(&install, "caroot", "", false, "Install the local CA in the system trust store.")
	mkcertCmd.Flags().BoolVarP(&install, "install", "", false, "Install the local CA in the system trust store.")
	mkcertCmd.Flags().BoolVarP(&uninstall, "uninstall", "", false, "Uninstall the local CA (but do not delete it).")
}
