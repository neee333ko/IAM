package config

import "github.com/neee333ko/IAM/internal/apiserver/option"

type Config struct {
	*option.Option
}

func NewConfig(option *option.Option) *Config {
	return &Config{
		Option: option,
	}
}
