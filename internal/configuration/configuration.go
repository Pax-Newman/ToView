package configuration

import (
	"os"
	"path/filepath"

	"github.com/Pax-Newman/toview/internal/filehelpers"
	"github.com/spf13/viper"
)

// Load a viper from a config file
func LoadConfig(path string) (*viper.Viper, error) {
	config := viper.New()

	homedir, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}

	// build the paths to the config files
	configdir := filepath.Join(homedir, "/.config/toview/")
	configPath := filepath.Join(configdir, path)

	filename := filehelpers.GetFilename(path)
	filetype, err := filehelpers.GetExtension(path)
	if err != nil {
		return nil, err
	}

	config.SetConfigName(filename)
	config.SetConfigType(filetype)
	// add config file path for ~/.config/toview/
	config.AddConfigPath(configPath)
	config.AddConfigPath("./.config/")
	config.AddConfigPath(".")
	err = config.ReadInConfig()
	if err != nil {
		return config, err
	}

	return config, nil
}

// Unmarshals from a path into a generic struct which it then returns
func UnmarshalFromPath[T interface{}](path string) (T, error) {
	var newType T

	config, err := LoadConfig(path)
	if err != nil {
		return newType, err
	}

	err = config.Unmarshal(&newType)
	if err != nil {
		return newType, err
	}

	return newType, nil
}
