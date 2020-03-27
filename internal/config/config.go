package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	NatsUrl            string `yaml:"natsUrl"`
	NatsCreatedSubject string `yaml:"natsCreatedSubject"`
	NatsQueueName      string `yaml:"natsQueueName"`
	RollbarToken       string `yaml:"rollbarToken"`
}

func (c *Config) PrettyPrint() string {
	configStr := fmt.Sprintf(
		"access_key: <hidden> \nurl: %s \ntests: %s %s", c.NatsUrl, c.NatsQueueName, c.NatsCreatedSubject)
	return configStr
}

func GetConf() *Config {
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error(err)
	}
	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		logrus.Errorf("unable to decode into config struct, %v", err)
	}
	return conf
}
