package config

import "github.com/spf13/viper"

type Config struct {
	NatsUrl            string
	NatsCreatedSubject string
	NatsQueueName      string
	RollbarToken       string
}

func GetConf() *Config {
	viper.AutomaticEnv()
	conf := &Config{
		NatsUrl:            viper.GetString("nats_url"),
		NatsCreatedSubject: viper.GetString("nats_created_subject"),
		NatsQueueName:      viper.GetString("nats_queue_name"),
		RollbarToken:       viper.GetString("rollbar_token"),
	}
	return conf
}
