package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type cmdFn func(cmd *cobra.Command, args []string) error

func runCmd(fn cmdFn) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := fn(cmd, args); err != nil {
			fmt.Printf("Error: %s\n\n", err.Error())

			_ = cmd.Help()

			os.Exit(1)
		}
	}
}
