package config

import "github.com/neee333ko/IAM/internal/authzserver/option"

type Config struct {
	*option.Option
}

func CreateConfigFromOption(option *option.Option) *Config {
	return &Config{Option: option}
}
