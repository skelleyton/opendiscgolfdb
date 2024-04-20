package utils

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	DbUser  string `dotenv:"DB_USER"`
	DbPass  string `dotenv:"DB_PASSWORD"`
	ConnStr string `dotenv:"CONN_STR"`
}

type DotenvConfig struct {
	Config *Config
}

func NewDotenvConfig(path string) *DotenvConfig {
	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	var fullPath string

	if path != "" {
		if !strings.HasPrefix(path, "/") {
			fullPath = fmt.Sprintf("%s/%s", dir, path)
		} else {
			fullPath = path
		}
	} else {
		fullPath = fmt.Sprintf("%s/.env", dir)
	}

	fileByte, err := os.ReadFile(fullPath)

	if err != nil {
		log.Print(err)
		log.Fatalf("Failed to read dotnev file %s", fullPath)
	}

	fileString := string(fileByte)
	configMap := make(map[string]string)

	for _, val := range strings.Split(fileString, "\n") {
		splitVal := strings.Split(val, "=")
		configMap[splitVal[0]] = splitVal[1]
	}

	config := insertToStruct(configMap)

	return &DotenvConfig{config}
}

func insertToStruct(configMap map[string]string) *Config {
	config := &Config{}

	structValue := reflect.ValueOf(config).Elem()
	configType := structValue.Type()

	structFields := reflect.VisibleFields(configType)

	for key, configVal := range configMap {
		for _, val := range structFields {
			dotenvTag := val.Tag.Get("dotenv")
			if dotenvTag == key {
				structField := structValue.FieldByName(val.Name)
				valType := reflect.ValueOf(configVal)

				if structField.IsValid() && structField.CanSet() && structField.Type() == valType.Type() {
					structField.Set(valType)
				}
			}
		}
	}

	return config
}
