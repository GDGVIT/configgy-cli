package login

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-cli/utils"
	"github.com/GDGVIT/configgy-cli/utils/client"
	"github.com/GDGVIT/configgy-cli/utils/config"
	"github.com/GDGVIT/configgy-cli/utils/crypto"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login to configgy",
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
			secretPath, err := cmd.Flags().GetString("secretpath")
			if err != nil {
				return err
			}

			if login == "" || password == "" {
				fmt.Println("Welcome to configgy. Please enter your details to login.")
				fmt.Printf("Enter your email: ")
				fmt.Scanln(&login)
				fmt.Printf("Enter your password: ")
				fmt.Scanln(&password)
			}

			// check if secret key exists
			if secretPath == "" {
				secretPath = filepath.Join(config.GetConfiggyCliConfigDirectory(), "keys", "secret.key")
				if !utils.FilePathExists(secretPath) {
					fmt.Println("Secret key not found. Please provide the path to the secret key file.")
					return nil
				}
			}

			// read secret key file for salt
			salt, err := os.ReadFile(secretPath)
			if err != nil {
				return err
			}

			key, _, err := crypto.GeneratePassword(login+password, salt)

			// api call to login
			loginRequest := api.LoginJSONRequestBody{
				Email:    types.Email(login),
				Password: key,
			}
			loginResponse := api.LoginResponse{}

			// create a resty client
			c, err := client.CreateRestyClient("", "")
			if err != nil {
				return err
			}
			resp, err := c.R().SetBody(loginRequest).SetResult(&loginResponse).Post("/user/login")
			if resp.IsError() {
				return fmt.Errorf("login failed: %s", resp.Status())
			}
			if err != nil {
				return err
			}

			// save the token to config file
			config.SetCredentials("auth", "access_token", loginResponse.Token.AccessToken)
			config.SetCredentials("auth", "user_id", *loginResponse.Token.UserId)

			fmt.Println("Login successful.")
			return nil
		},
	}
}

func init() {
	// add commands here with RootCmd().AddCommand()
	RootCmd().Flags().StringP("email", "e", "", "Email to login")
	RootCmd().Flags().StringP("password", "p", "", "Password to login")
	RootCmd().Flags().StringP("secretpath", "s", "", "Path to seecret key file")
}

func Execute() {
	if err := RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
