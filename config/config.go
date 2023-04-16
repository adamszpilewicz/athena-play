package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type AWSConfig struct {
	Key    string `mapstructure:"aws_key"`
	Secret string `mapstructure:"aws_secret"`
	Region string `mapstructure:"aws_region"`
	Bucket string `mapstructure:"aws_bucket"`
}

type Config struct {
	AWS AWSConfig `mapstructure:"aws_config"`
}

// Validate checks if the AWSConfig contains all required fields.
func (c AWSConfig) Validate() error {
	if c.Key == "" {
		return errors.New("missing aws_key")
	}
	if c.Secret == "" {
		return errors.New("missing aws_secret")
	}
	if c.Region == "" {
		return errors.New("missing aws_region")
	}
	return nil
}

func LoadConfig(configPath string) (Config, error) {
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("failed to read the config file: %v", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal the config: %v", err)
	}

	// Validate the AWSConfig
	err = config.AWS.Validate()
	if err != nil {
		return config, fmt.Errorf("invalid AWS configuration: %v", err)
	}

	return config, nil
}
