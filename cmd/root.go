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

var configFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(
		&configFile, "config", "", "config file (default is config.yaml)")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/sailor")
		viper.AddConfigPath("$HOME/.sailor")
	}
	viper.SetEnvPrefix("boatswain")
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("unable to read config: %v\n", err)
		os.Exit(1)
	}
	conf := config.GetConf()
	logrus.Infof("running with config: %v", conf)

}
