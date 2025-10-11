package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/neee333ko/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const configName = "Config"

var config string

// AddConfigFlag add config flag to fs and set cobra.OnInitialize to read in config when execute
func AddConfigFlag(fs *pflag.FlagSet, basename string) {
	fs.StringVarP(&config, configName, "c", config,
		fmt.Sprintf("Config path for %s, supports YAML only.",
			FormatName(basename)))

	cobra.OnInitialize(func() {
		if config != "" {
			viper.SetConfigFile(config)
		} else {
			viper.SetConfigName(FormatName(basename))
			viper.SetConfigType("yaml")

			viper.AddConfigPath(".")

			home, _ := os.UserHomeDir()
			viper.AddConfigPath(filepath.Join(home, strings.Split(basename, "-")[0]))
			viper.AddConfigPath(filepath.Join("/etc", strings.Split(basename, "-")[0]))
		}

		if err := viper.ReadInConfig(); err != nil {
			log.Warnf("%s\n", color.YellowString("WARNING: ", err.Error()))
		}
	})
}

func printConfigFlag() string {
	return fmt.Sprintf("Config Path: %s", config)
}
