package main

import (
	"ctRestClient/app"
	"ctRestClient/config"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/term"
)

func main() {
	var configFilePath string
	var outputDirectory string
	var keepassDbFilePath string

	flag.StringVar(&configFilePath, "c", "config.yml", "the config file path")
	flag.StringVar(&outputDirectory, "o", getExecutablePath(), "the output directory")
	flag.StringVar(&keepassDbFilePath, "k", "passwords.kdbx", "the Keepass DB file path")
	flag.Parse()

	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load config from path %s: %v", configFilePath, err)
	}

	keepassDbPassword, err := getPasswordFromUser()
	if err != nil {
		log.Fatalf("Failed to get password: %v", err)
	}

	appLogger := app.NewLogger()
	err = app.NewInstancesProcessor(*config, outputDirectory, appLogger).Process(
		app.NewGroupExporter(),
		app.NewCSVFileWriter(),
		app.NewKeepassCli(keepassDbFilePath, keepassDbPassword, appLogger),
	)
	if err != nil {
		log.Fatalf("Failed to process instances: %v", err)
	}
}

func getExecutablePath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	return filepath.Dir(exePath)
}

func getPasswordFromUser() (string, error) {
	fmt.Print("Enter Keepass database password: ")

	password, err := term.ReadPassword(0) // 0 is stdin file descriptor
	if err != nil {
		return "", fmt.Errorf("failed to read password: %v", err)
	}
	fmt.Println()

	return string(password), nil
}
