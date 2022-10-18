/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package runtime

import (
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
	"liferay.com/lcectl/flags"
	lio "liferay.com/lcectl/io"
	"liferay.com/lcectl/mkcert"
	"liferay.com/lcectl/prereq"
)

// runtimeCmd represents the runtime command
var runtimeCmd = &cobra.Command{
	Use:     "runtime COMMAND [OPTIONS] [ARG...]",
	Short:   "Operations to control the runtime environment",
	Aliases: []string{"rt"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		prereq.Prereq(flags.Verbose)

		if cmd.Name() != "mkcert" {
			lfrdevCrtFile := path.Join(viper.GetString(constants.Const.RepoDir), "/k8s/tls/lfrdev.crt")
			lfrdevKeyFile := path.Join(viper.GetString(constants.Const.RepoDir), "/k8s/tls/lfrdev.key")
			lfrdevRootCAFile1 := path.Join(viper.GetString(constants.Const.RepoDir), "/k8s/tls/", mkcert.GetRootName())
			lfrdevRootCAFile2 := path.Join(viper.GetString(constants.Const.RepoDir), "/k8s/tls/", mkcert.GetRootName())
			lfrdevRootCAFile3 := path.Join(viper.GetString(constants.Const.RepoDir), "/k8s/tls/", mkcert.GetRootName())

			if !lio.Exists(lfrdevCrtFile) || !lio.Exists(lfrdevKeyFile) || !lio.Exists(lfrdevRootCAFile1) || !lio.Exists(lfrdevRootCAFile2) || !lio.Exists(lfrdevRootCAFile3) {
				log.Fatalf("Missing one or more local certificates.  Execute 'runtime mkcert' command to generate one.")
			}

			crt, key, err := mkcert.LoadX509KeyPair(lfrdevCrtFile, lfrdevKeyFile)
			if crt == nil || key == nil || err != nil {
				log.Fatalf("Could not load x509 key pair: %s", err)
			}
			lrdevDomain := viper.GetString(constants.Const.TlsLfrdevDomain)
			if crt.DNSNames[0] != "*."+lrdevDomain {
				log.Fatalf("Generated certificate DNSName does not match configured domain: %s != %s\nPlease run 'runtime mkcert' command again.", crt.DNSNames[0], lrdevDomain)
			}
		}
	},
}

func AddRuntimeCmd(cmd *cobra.Command) {
	cmd.AddCommand(runtimeCmd)
}
