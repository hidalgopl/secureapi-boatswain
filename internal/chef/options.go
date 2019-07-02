package chef

import "github.com/kelseyhightower/envconfig"

//TODO(core) use CHEF_SCHEDULER env prefix in chart/deployment
const envPrefix = ""

type Options struct {
	Namespace        string `envconfig:"CHEF_NAMESPACE" default:"default"`
	ChefImage        string `envconfig:"CHEF_IMAGE" required:"true"`
	ChefTag          string `envconfig:"CHEF_TAG" required:"true"`
	BrowserSyncImage string `envconfig:"BROWSER_SYNC_IMAGE" required:"true"`
	LocalMode        bool   `envconfig:"LOCAL_MODE" default:"false"`
	ChefDryads       string `envconfig:"CHEF_DRYADS" required:"true"`

	HeadlessChefImage          string `envconfig:"HEADLESS_CHEF_IMAGE" required:"true"`
	HeadlessChefTag            string `envconfig:"HEADLESS_CHEF_TAG" required:"true"`
	HeadlessChefServiceAccount string `envconfig:"HEADLESS_CHEF_SERVICE_ACCOUNT" required:"true"`
}

func OptionsFromEnv() (Options, error) {
	var options Options
	err := envconfig.Process(envPrefix, &options)
	return options, err
}

func PrintEnvOptionsUsage() error {
	return envconfig.Usage(envPrefix, &Options{})
}
