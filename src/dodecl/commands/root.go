package commands

import "github.com/spf13/cobra"

// RootCmd is root command for dodecl.
var RootCmd = &cobra.Command{
	Use:   "dodecl",
	Short: "dodecl is a declarative state tool for DigitalOcean",
	Long:  "dodecl manages multiple DigitalOcean resources at once.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
}
