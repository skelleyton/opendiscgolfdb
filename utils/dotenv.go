package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type DotenvConfig struct {
	Config map[string]string
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
	
	configArray := strings.Split(fileString, "\n")
	
	config := make(map[string]string)
	
	for _, val := range configArray {
		splitVal := strings.Split(val, "=")
		config[splitVal[0]] = splitVal[1]
	}
	
	return &DotenvConfig{config}
}