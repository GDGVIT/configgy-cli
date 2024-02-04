package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "credential",
		Short: "Manage credentials",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from credential.")
		},
	}
}

func init() {
	// add commands here with RootCmd().AddCommand()
}
