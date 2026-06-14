package app

import (
	"ctRestClient/internal/config"
	"ctRestClient/internal/csv"
	"ctRestClient/internal/dataprovider"
	"ctRestClient/internal/httpclient"
	"ctRestClient/internal/logger"
	"ctRestClient/internal/rest"
	"fmt"
	"os"
	"path/filepath"
)

type exportGroupTask struct {
	groupExporter          GroupExporter
	csvWriter              csv.CSVFileWriter
	rootDir                string
	fileDataProvider       dataprovider.FileDataProvider
	blocklistsDataProvider dataprovider.BlockListDataProvider
	logger                 logger.Logger
}

func NewExportGroupTask(
	groupExporter GroupExporter,
	csvWriter csv.CSVFileWriter,
	rootDir string,
	fileDataProvider dataprovider.FileDataProvider,
	blocklistsDataProvider dataprovider.BlockListDataProvider,
	logger logger.Logger,
) exportGroupTask {
	return exportGroupTask{
		groupExporter:          groupExporter,
		csvWriter:              csvWriter,
		rootDir:                rootDir,
		fileDataProvider:       fileDataProvider,
		blocklistsDataProvider: blocklistsDataProvider,
		logger:                 logger,
	}
}

// Implement InstanceTask.Execute method
func (p exportGroupTask) Execute(instance config.Instance, httpClient httpclient.HTTPClient) error {
	groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
	groupEndpoint := rest.NewGroupEndpoint(httpClient)
	dynamicGroupsEndpoint := rest.NewDynamicGroupsEndpoint(httpClient)
	personEndpoint := rest.NewPersonsEndpoint(httpClient)

	for _, group := range instance.Groups {
		p.logger.Info("")
		p.logger.Info(fmt.Sprintf("  processing group '%s'", group.Name))

		persons, err := p.groupExporter.ExportGroupMembers(
			group.Name,
			groupsEndpoint,
			groupEndpoint,
			dynamicGroupsEndpoint,
			personEndpoint,
		)
		if err != nil {
			if _, ok := err.(*GroupNotActiveError); ok {
				p.logger.Warn("      skipping csv creation since the group is not active")
				continue
			} else {
				p.logger.Error(fmt.Sprintf("      failed to get person information: %v", err))
				continue
			}
		}

		if len(persons) == 0 {
			p.logger.Info("      the group is empty")
			continue
		} else {
			p.logger.Info(fmt.Sprintf("      the group has %d persons", len(persons)))
		}

		if p.blocklistsDataProvider.BlockListExists(group) {
			p.logger.Info(fmt.Sprintf("      using blocklist '%s'", group.BlocklistFileName()))
		}

		personData, err := csv.NewPersonData(persons, group, p.fileDataProvider, p.blocklistsDataProvider, p.logger)
		if err != nil {
			p.logger.Error(fmt.Sprintf("      failed to extract persons: %v", err))
			continue
		}

		err = os.MkdirAll(filepath.Join(p.rootDir, instance.Hostname), 0755)
		if err != nil {
			p.logger.Error(fmt.Sprintf("     failed to create directory: %v", err))
			continue
		}

		csvFilePath := filepath.Join(
			p.rootDir,
			instance.Hostname,
			group.CSVFileName(),
		)

		err = p.csvWriter.Write(csvFilePath, personData.Header(), personData.Records())
		if err != nil {
			p.logger.Error(fmt.Sprintf("    failed to write csv file: %v", err))
			continue
		}
	}

	return nil
}
