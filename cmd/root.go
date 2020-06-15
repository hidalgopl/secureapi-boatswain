package cmd

import (
	"github.com/hidalgopl/secureapi-boatswain/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/rollbar/rollbar-go"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "boatswain",
	Short: "Boatswain is test suite worker for SecureAPI",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func Execute() {
	conf := config.GetConf()
	rollbar.SetToken(conf.RollbarToken)
	rollbar.SetEnvironment("production")
	//rollbar.SetCodeVersion("")
	rollbar.WrapAndWait(main)

}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetEnvPrefix("boatswain")
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	conf := config.GetConf()
	logrus.Infof("running with config: %v", conf)

}
