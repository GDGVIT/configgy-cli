package cmd

import (
	"fmt"

	"github.com/GDGVIT/configgy-cli/authn/pkg/login"
	"github.com/GDGVIT/configgy-cli/authn/pkg/signup"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Commands for authentication",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from auth.")
		},
	}
}

func init() {
	// add commands here with RootCmd().AddCommand()
	RootCmd().AddCommand(login.RootCmd())
	RootCmd().AddCommand(signup.RootCmd())
}

func Execute() {
	if err := RootCmd().Execute(); err != nil {
		panic(err)
	}
}
