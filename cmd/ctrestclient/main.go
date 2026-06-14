package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"ctRestClient/internal/app"
	"ctRestClient/internal/config"
	"ctRestClient/internal/csv"
	"ctRestClient/internal/dataprovider"
	"ctRestClient/internal/logger"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	configFilePath     string
	dataDir            string
	outputDir          string
	keepassDbFilePath  string
	rootCmd            = &cobra.Command{
		Use:   "ctRestClient",
		Short: "A CLI tool for exporting data using Keepass and CSV",
		Long:  `ctRestClient is a command-line application that exports data from configured instances using Keepass for secure credential management and CSV for output.`,
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "config.yml", "Path to the config file")
	rootCmd.PersistentFlags().StringVarP(&dataDir, "data-dir", "d", getDefaultDataDir(), "Path to the data directory")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "o", getDefaultOutputDir(), "Path to the output directory")
	rootCmd.PersistentFlags().StringVarP(&keepassDbFilePath, "keepass", "k", "passwords.kdbx", "Path to the Keepass database file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}

func run() {
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

	exportGroupTask := app.NewExportGroupTask(
		app.NewGroupExporter(),
		csv.NewCSVFileWriter(),
		rootDir,
		dataprovider.NewFileDataProvider(filepath.Join(dataDir, "mappings/persons")),
		dataprovider.NewBlockListDataProvider(filepath.Join(dataDir, "blocklists"), appLogger),
		appLogger,
	)

	err = app.NewInstancesProcessor(
		config.Instances,
		keepassCli,
		appLogger,
	).Process(
		exportGroupTask,
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
	if password, exists := os.LookupEnv("KEY_PASS_PASSWORD"); exists {
		return password, nil
	}

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
