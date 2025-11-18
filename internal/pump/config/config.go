package config

import "github.com/neee333ko/IAM/internal/pump/option"

type Config struct {
	*option.Option
}

func CreateConfigFromOption(option *option.Option) (*Config, error) {
	return &Config{Option: option}, nil
}
