package main

import (
	"ctRestClient/app"
	"ctRestClient/config"
	"ctRestClient/csv"
	"ctRestClient/data_provider"
	"ctRestClient/logger"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

func main() {
	var configFilePath string
	var dataDir string
	var outputDir string
	var keepassDbFilePath string

	flag.StringVar(&configFilePath, "c", "config.yml", "the config file path")
	flag.StringVar(&dataDir, "d", getDefaultDataDir(), "the data directory")
	flag.StringVar(&outputDir, "o", getDefaultOutputDir(), "the output directory")
	flag.StringVar(&keepassDbFilePath, "k", "passwords.kdbx", "the Keepass DB file path")
	flag.Parse()

	rootDir := filepath.Join(outputDir, time.Now().Format("2006.01.02_15-04-05"))
	err := os.MkdirAll(rootDir, 0755)
	if err != nil {
		log.Fatalf("    failed to create directory: %v", err)
	}

	logFile := filepath.Join(rootDir, "ctRestClient.log")

	appLogger := logger.NewLogger(logFile)
	logGeneralInfo(appLogger, getCurrentUserName(), getCurrentOSName(), getDate())

	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		appLogger.Fatal(fmt.Sprintf("Failed to load config from path %s: %v", configFilePath, err))
	}

	keepassDbPassword, err := getPasswordFromUser()
	if err != nil {
		appLogger.Fatal(fmt.Sprintf("Failed to get password: %v", err))
	}

	keepassCli, err := app.NewKeepassCli(keepassDbFilePath, keepassDbPassword, appLogger)
	if err != nil {
		appLogger.Fatal(fmt.Sprintf("Failed to initialize Keepass CLI: %v", err))
	}

	validPassword, err := keepassCli.IsPasswordValid(keepassDbPassword)
	if err != nil {
		appLogger.Fatal(fmt.Sprintf("Failed check keepass password: %v", err))
	}

	if !validPassword {
		appLogger.Fatal("The keepass password is invalid")
	}

	err = app.NewInstancesProcessor(
		*config,
		appLogger,
	).Process(
		app.NewGroupExporter(),
		csv.NewCSVFileWriter(),
		rootDir,
		data_provider.NewFileDataProvider(filepath.Join(dataDir, "mappings/persons")),
		data_provider.NewBlockListDataProvider(filepath.Join(dataDir, "blocklists"), appLogger),
		keepassCli,
	)
	if err != nil {
		appLogger.Fatal(fmt.Sprintf("Failed to process instances: %v", err))
	}
}

func getDefaultOutputDir() string {
	executableDir := getExecutableDir()
	return filepath.Join(executableDir, "..", "exports")
}

func getDefaultDataDir() string {
	executableDir := getExecutableDir()
	return filepath.Join(executableDir, "..", "data")
}

func getExecutableDir() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	executabelDir := filepath.Dir(exePath)
	return executabelDir
}

func getPasswordFromUser() (string, error) {
	fmt.Print("Enter Keepass database password: ")

	// Use the appropriate file descriptor based on the platform
	fd := int(syscall.Stdin)
	password, err := term.ReadPassword(fd)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %v", err)
	}
	fmt.Println()

	return string(password), nil
}

func getCurrentUserName() string {
	currentUser, err := user.Current()
	if err != nil {
		return "unknown"
	}
	return currentUser.Username
}

func getCurrentOSName() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return fmt.Sprintf("Unknown (%s)", runtime.GOOS)
	}
}

func getDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func logGeneralInfo(logger logger.Logger, user string, os string, date string) {
	boxLength := 70
	userInfo := fmt.Sprintf("User: '%s'", user)
	userInfoLength := len(userInfo)

	osInfo := fmt.Sprintf("OS:   '%s'", os)
	osInfoLength := len(osInfo)

	dateInfo := fmt.Sprintf("Date: '%s'", date)
	dateInfoLength := len(dateInfo)

	border := strings.Repeat("-", boxLength)

	logger.Info("")
	logger.Info(fmt.Sprintf("+%s+", border))
	logger.Info(fmt.Sprintf("| %s "+strings.Repeat(" ", boxLength-userInfoLength-2)+"|", userInfo))
	logger.Info(fmt.Sprintf("| %s "+strings.Repeat(" ", boxLength-osInfoLength-2)+"|", osInfo))
	logger.Info(fmt.Sprintf("| %s "+strings.Repeat(" ", boxLength-dateInfoLength-2)+"|", dateInfo))
	logger.Info(fmt.Sprintf("+%s+", border))
	logger.Info("")
}
