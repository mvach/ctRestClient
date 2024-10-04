package app

import (
	"ctRestClient/config"
	"ctRestClient/csv"
	"ctRestClient/httpclient"
	"ctRestClient/logger"
	"ctRestClient/rest"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type InstancesProcessor interface {
	Process(groupExporter GroupExporter, csvWriter csv.CSVFileWriter, keepassCli KeepassCli) error
}

type instancesProcessor struct {
	config          config.Config
	outputDirectory string
	logger          logger.Logger
}

func NewInstancesProcessor(
	config config.Config,
	outputDirectory string,
	logger logger.Logger,
) InstancesProcessor {
	return instancesProcessor{
		config:          config,
		outputDirectory: outputDirectory,
		logger:          logger,
	}
}

func (p instancesProcessor) Process(groupExporter GroupExporter, csvWriter csv.CSVFileWriter, keepassCli KeepassCli) error {
	// define the root dir location via parameter that is by default beside the executable
	rootDir := filepath.Join(p.outputDirectory, time.Now().Format("2006.01.02_15-04-05"))

	for _, instance := range p.config.Instances {
		p.logger.Info(fmt.Sprintf("processing instance '%s'", instance.Hostname))

		token, err := keepassCli.GetPassword(instance.TokenName)
		if err != nil {
			p.logger.Warn(fmt.Sprintf("  skipping export, failed to get token with name '%s' from Keepass. Err: %v", instance.TokenName, err))
			continue
		}

		httpClient := httpclient.NewHTTPClient(instance.Hostname, token)
		groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
		personEndpoint := rest.NewPersonsEndpoint(httpClient)

		for _, group := range instance.Groups {
			p.logger.Info(fmt.Sprintf("  processing group '%s'", group.Name))

			err := os.MkdirAll(filepath.Join(rootDir, instance.Hostname), 0755)
			if err != nil {
				p.logger.Error(fmt.Sprintf("    failed to create directory: %v", err))
				continue
			}

			persons, err := groupExporter.ExportGroupMembers(
				group.Name,
				groupsEndpoint,
				personEndpoint,
			)
			if err != nil {
				p.logger.Error(fmt.Sprintf("    failed to get person information: %v", err))
				continue
			}

			if len(persons) == 0 {
				p.logger.Info("    the group is empty")
				continue
			} else {
				p.logger.Info(fmt.Sprintf("    the group has %d persons", len(persons)))
			}

			personData, err := csv.NewPersonData(persons, group.Fields, p.logger)
			if err != nil {
				p.logger.Error(fmt.Sprintf("    failed to extract persons: %v", err))
				continue
			}

			csvFilePath := filepath.Join(
				rootDir,
				instance.Hostname,
				group.SanitizedGroupName()+".csv",
			)

			err = csvWriter.Write(csvFilePath, personData.Header(), personData.Records())
			if err != nil {
				p.logger.Error(fmt.Sprintf("    failed to write csv file: %v", err))
				continue
			}
		}
	}

	return nil
}
