package app

import (
	"github.com/spf13/cobra"
)

func helpCommand(basename string) *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "help",
		Long: `help command for ` + FormatName(basename) + ` ,just type ` + FormatName(basename) +
			` help [command] to get more information about [command]`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			arg := []string{args[0]}
			sc, _, err := cmd.Root().Find(arg)

			if sc != nil && err == nil {
				sc.Usage()
			} else {
				cmd.Root().Usage()
			}
		},
	}
}
