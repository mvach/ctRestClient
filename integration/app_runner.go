package integration

import (
	"ctRestClient/app"
	"ctRestClient/config"
	"ctRestClient/data_provider"
	"ctRestClient/internal/csv"
	"ctRestClient/internal/logger"
	"path/filepath"
)

// RunApplicationWrapper wraps the main application logic for integration testing
func RunApplicationWrapper(config *config.Config, rootDir string, dataDir string, keepassDbFilePath string, keepassDbPassword string, appLogger logger.Logger) error {
	keepassCli, err := app.NewKeepassCli(keepassDbFilePath, keepassDbPassword, appLogger)
	if err != nil {
		return err
	}
	return app.NewInstancesProcessor(
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
}
