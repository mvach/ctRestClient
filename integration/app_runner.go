package integration

import (
	"ctRestClient/app"
	"ctRestClient/config"
	"ctRestClient/csv"
	"ctRestClient/logger"
	"path/filepath"
)

// RunApplicationWrapper wraps the main application logic for integration testing
func RunApplicationWrapper(config *config.Config, rootDir string, dataDir string, keepassDbFilePath string, keepassDbPassword string, appLogger logger.Logger) error {
	return app.NewInstancesProcessor(
		*config,
		appLogger,
	).Process(
		app.NewGroupExporter(),
		csv.NewCSVFileWriter(),
		rootDir,
		csv.NewFileDataProvider(filepath.Join(dataDir, "persons")),
		app.NewKeepassCli(keepassDbFilePath, keepassDbPassword, appLogger),
	)
}
