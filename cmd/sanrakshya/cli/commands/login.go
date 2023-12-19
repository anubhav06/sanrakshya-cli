package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/anchore/clio"
	"github.com/anubhav06/sanrakshya-cli/internal"
)

const (
	loginExample = `
  Login to the Sanrakshya Web App. Stores the credentials in the environment variables:
  sanrakshya login --username <username> --password <password>
  `
)

//nolint:dupl
func Login(app clio.Application) *cobra.Command {
	id := app.ID()

	opts := defaultSubmitOptions()

	var username, password string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Sanrakshya Enterprise Server",
		Long:  "Login to Sanrakshya Web App",
		Example: internal.Tprintf(loginExample, map[string]interface{}{
			"appName": id.Name,
			"command": "submit",
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			// restoreStdout := ui.CaptureStdoutToTraceLog()
			// defer restoreStdout()

			return runLogin(username, password)
		},
	}
	cmd.Flags().StringVarP(&username, "username", "u", "", "Your username")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Your password")

	return app.SetupCommand(cmd, opts)

}

func runLogin(username string, password string) error {

	fmt.Println("Logging in...")

	// Create a map with the username and password
	data := map[string]string{"account_id": username, "password": password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Make API call to get ACCOUNT_ID and SECRET_KEY
	resp, err := http.Post("http://127.0.0.1:8000/api/get-key/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the API call was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API call failed with status code: %d", resp.StatusCode)
	}

	// Parse the response body to get ACCOUNT_ID and SECRET_KEY
	var keys struct {
		AccountID string `json:"account_id"`
		SecretKey string `json:"secret_key"`
	}
	err = json.NewDecoder(resp.Body).Decode(&keys)
	if err != nil {
		return err
	}

	fmt.Println(keys.AccountID, keys.SecretKey)

	// Get the installation location of the CLI tool
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	installDir := filepath.Dir(exe)

	// Create a file to store the credentials
	credsFile := filepath.Join(installDir, "creds.json")
	creds := struct {
		AccountID string `json:"account_id"`
		SecretKey string `json:"secret_key"`
	}{
		AccountID: keys.AccountID,
		SecretKey: keys.SecretKey,
	}
	credsBytes, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	// Write the credentials to the file
	err = ioutil.WriteFile(credsFile, credsBytes, 0600)
	if err != nil {
		return err
	}

	fmt.Println("Login successful")

	return nil
}
