package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	NatsUsername       string
	NatsUrl            string
	NatsPass           string
	NatsCreatedSubject string
	NatsQueueName      string
}

func (c *Config) PrettyPrint() string {
	configStr := fmt.Sprintf(
		"username: %s \naccess_key: <hidden> \nurl: %s \ntests: %s", c.NatsUsername, c.NatsUrl, c.NatsQueueName, c.NatsCreatedSubject)
	return configStr
}

func GetConf() *Config {
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%v", err)
	}
	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}
	return conf
}
