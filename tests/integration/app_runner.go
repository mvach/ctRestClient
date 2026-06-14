package integration

import (
	"ctRestClient/internal/app"
	"ctRestClient/internal/config"
	"ctRestClient/internal/csv"
	"ctRestClient/internal/dataprovider"
	"ctRestClient/internal/logger"
	"path/filepath"
)

// RunApplicationWrapper wraps the main application logic for integration testing
func RunApplicationWrapper(config *config.Config, rootDir string, dataDir string, keepassDbFilePath string, keepassDbPassword string, appLogger logger.Logger) error {
	keepassCli, err := app.NewKeepassCli(keepassDbFilePath, keepassDbPassword, appLogger)
	if err != nil {
		return err
	}

	exportGroupTask := app.NewExportGroupTask(
		app.NewGroupExporter(),
		csv.NewCSVFileWriter(),
		rootDir,
		dataprovider.NewFileDataProvider(filepath.Join(dataDir, "mappings/persons")),
		dataprovider.NewBlockListDataProvider(filepath.Join(dataDir, "blocklists"), appLogger),
		appLogger,
	)

	return app.NewInstancesProcessor(
		config.Instances,
		keepassCli,
		appLogger,
	).Process(
		exportGroupTask,
	)
}
