package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Parse knows how to parse file to config struct
func Parse(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(v); err != nil {
		return err
	}

	return nil
}
