package cmd

import (
	"github.com/hidalgopl/secureapi-boatswain/internal/config"
	"github.com/hidalgopl/secureapi-boatswain/internal/listener"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use: "run",
	Short: "Runs boatswain service",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConf()
		err := listener.Listen(conf)
		if err != nil {
			logrus.Error(err)
		}
	},
}