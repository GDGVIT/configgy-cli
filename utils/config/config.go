package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var configName = "config"
var fileName = configName + ".toml"
var directoryName = ".configgy-cli"

func init() {
	// Set the name of the config file
	viper.SetConfigName(configName)

	// Set the path to the directory containing the config file
	viper.AddConfigPath(filepath.Join(GetConfiggyCliConfigDirectory()))

	// Create the directory and config file if they don't exist
	if err := ensureConfigFile(); err != nil {
		panic(fmt.Errorf("failed to initialize config file: %s", err))
	}

	// Read in the config file
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}
}

func ensureConfigFile() error {
	// Get the path to the config file
	configPath := filepath.Join(GetConfiggyCliConfigDirectory(), fileName)

	// Check if the config file exists
	if !filePathExists(configPath) {
		// Config file does not exist, create the directory and file
		err := os.MkdirAll(filepath.Dir(configPath), 0600)
		if err != nil {
			return err
		}

		_, err = os.Create(configPath)
		if err != nil {
			return err
		}

		// check for "keys" directory in config, if it does not exist, create one and copy the key there
		keysFolderPath := filepath.Join(GetConfiggyCliConfigDirectory(), "keys")
		if !filePathExists(keysFolderPath) {
			// create the directory
			err := os.MkdirAll(keysFolderPath, 0600)
			if err != nil {
				panic(fmt.Errorf("failed to create keys directory: %s", err))
			}
		}

		fmt.Println("Created config file:", configPath)
	}
	return nil
}

func SetCredentials(appName string, key string, value string) {
	appConfig := viper.GetStringMap(appName)
	if appConfig == nil {
		appConfig = make(map[string]interface{})
	}

	if _, ok := appConfig[key]; ok {
		// Key already exists, replace the value
		appConfig[key] = value
	} else {
		// Key doesn't exist, add a new key-value pair
		appConfig[key] = value
	}

	viper.Set(appName, appConfig)

	err := viper.WriteConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error writing config file: %s", err))
	}
}

func SetKeyPath(appName string, configKey string, keyPath string) {
	appConfig := viper.GetStringMap(appName)
	if appConfig == nil {
		appConfig = make(map[string]interface{})
	}

	// Check if the key path exists
	if !filePathExists(keyPath) {
		panic(fmt.Errorf("key path %s does not exist", keyPath))
	}
	sourceFile, err := os.Open(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to open key at: %s: %s", keyPath, err))
	}

	// check for "keys" directory in config, if it does not exist, create one and copy the key there
	keysFolderPath := filepath.Join(GetConfiggyCliConfigDirectory(), "keys")
	if !filePathExists(keysFolderPath) {
		// create the directory
		err := os.MkdirAll(keysFolderPath, 0600)
		if err != nil {
			panic(fmt.Errorf("failed to create keys directory: %s", err))
		}
	}
	//copy the key to the keys directory
	keyPath = filepath.Join(keysFolderPath, fmt.Sprintf("%s.key", appName))
	// check if the key already exists
	if filePathExists(keyPath) {
		// key already exists, remove it
		err := os.Remove(keyPath)
		if err != nil {
			panic(fmt.Errorf("failed to remove existing key: %s", err))
		}
	}
	// create the an empty file to hold the key
	destinationFile, err := os.Create(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to create destination file: %s", err))
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		panic(fmt.Errorf("failed to copy key: %s", err))
	}

	appConfig[configKey] = keyPath

	viper.Set(appName, appConfig)

	err = viper.WriteConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error writing config file: %s", err))
	}
}

func GetCredentials(appName string, key string) (string, error) {
	appConfig := viper.GetStringMapString(appName)
	if appConfig == nil {
		return "", fmt.Errorf("app %s not found", appName)
	}

	value, ok := appConfig[key]
	if !ok {
		return "", fmt.Errorf("key %s not found in app %s", key, appName)
	}

	return value, nil
}

func filePathExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist
			return false
		} // Some other error occurred
		panic(fmt.Errorf("failed to check if file exists: %s", err))
	}
	return true
}

func GetConfiggyCliConfigDirectory() string {
	return filepath.Join(homeDir(), directoryName)
}

// homeDir returns the path to the user's home directory.
func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}
