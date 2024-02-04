package signup

import (
	"fmt"
	"os"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-cli/utils/client"
	"github.com/GDGVIT/configgy-cli/utils/crypto"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "signup",
		Short: "Signup to configgy",
		RunE: func(cmd *cobra.Command, args []string) error {
			// business logic here
			login, err := cmd.Flags().GetString("email")
			if err != nil {
				return err
			}
			password, err := cmd.Flags().GetString("password")
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			if login == "" || password == "" || name == "" {
				fmt.Println("Welcome to configgy. Please enter your details to signup.")
				fmt.Printf("Enter your email: ")
				fmt.Scanln(&login)
				fmt.Printf("Enter your password: ")
				fmt.Scanln(&password)
				fmt.Printf("Enter your name: ")
				fmt.Scanln(&name)
			}

			key, skey, err := crypto.GeneratePassword(login+password, nil)
			if err != nil {
				return err
			}

			fmt.Println("Please save the following secret key safely. You will need it to login on another device.")
			fmt.Println(skey)

			pubKey, err := crypto.GenerateRSAKeyPair()
			if err != nil {
				return err
			}

			fmt.Println("Please save the following public key safely. You will need it to decrypt your data on another device.")
			fmt.Println(pubKey)

			// api call to signup
			c, err := client.CreateRestyClient("", "")
			if err != nil {
				return err
			}

			signupRequest := api.SignupJSONRequestBody{
				Email: types.Email(login),
				// key stretching with pbkdf2
				Password: key,
				Name:     name,
				// to be generated here with crypto
				PublicKey: pubKey,
			}

			signupResponse := api.GenericMessageResponse{}

			resp, err := c.R().SetResult(&signupResponse).SetBody(signupRequest).Post("/auth/signup")
			if resp.IsError() {
				return fmt.Errorf("Error: %v", resp.Error())
			}
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func init() {
	// add commands here with RootCmd().AddCommand()
	RootCmd().Flags().StringP("email", "e", "", "Email to signup")
	RootCmd().Flags().StringP("password", "p", "", "Password to signup")
	RootCmd().Flags().StringP("name", "n", "", "Name to signup")
}

func Execute() {
	if err := RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
