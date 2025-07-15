package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type (
	logConfig struct {
		Level       string `mapstructure:"Level"`
		Format      string `mapstructure:"Format"`
		ToFile      bool   `mapstructure:"ToFile"`
		Directory   string `mapstructure:"Directory"`
		Development bool   `mapstructure:"Development"`
	}
	natsConfig struct {
		Address    string   `mapstructure:"Address"`
		Username   string   `mapstructure:"Username"`
		Password   string   `mapstructure:"Password"`
		Subscribes []string `mapstructure:"Subscribes"`
	}
	weixinConfig struct {
		CorpID            string `mapstructure:"CorpId"`
		CorpSecret        string `mapstructure:"CorpSecret"`
		AgentID           int64  `mapstructure:"AgentId"`
		WebhookKey        string `mapstructure:"WebhookKey"`
		QYAPIHostOverride string `mapstructure:"QYAPIHostOverride"`
		TLSKeyLogFile     string `mapstructure:"TLSKeyLogFile"`
	}

	Config struct {
		Name   string        `mapstructure:"Name"`
		Log    *logConfig    `mapstructure:"Log"`
		Nats   *natsConfig   `mapstructure:"Nats"`
		Weixin *weixinConfig `mapstructure:"Weixin"`
	}
)

func getViper() *viper.Viper {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile("config.yaml")
	return v
}

func NewConfig() (*Config, error) {
	fmt.Println("Loading configuration")
	v := getViper()
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var config Config
	err = v.Unmarshal(&config)
	return &config, err
}
