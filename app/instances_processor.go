package app

import (
	"ctRestClient/config"
	"ctRestClient/csv"
	"ctRestClient/data_provider"
	"ctRestClient/httpclient"
	"ctRestClient/logger"
	"ctRestClient/rest"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type InstancesProcessor interface {
	Process(
		groupExporter GroupExporter,
		csvWriter csv.CSVFileWriter,
		rootDir string,
		personDataProvider data_provider.FileDataProvider,
		blocklistsDataProvider data_provider.BlockListDataProvider,
		keepassCli KeepassCli,
	) error
}

type instancesProcessor struct {
	config config.Config
	logger logger.Logger
}

func NewInstancesProcessor(
	config config.Config,
	logger logger.Logger,
) InstancesProcessor {
	return instancesProcessor{
		config: config,
		logger: logger,
	}
}

func (p instancesProcessor) Process(
	groupExporter GroupExporter,
	csvWriter csv.CSVFileWriter,
	rootDir string,
	fileDataProvider data_provider.FileDataProvider,
	blocklistsDataProvider data_provider.BlockListDataProvider,
	keepassCli KeepassCli,
) error {
	for _, instance := range p.config.Instances {

		p.logTitle(instance)

		token, err := keepassCli.GetPassword(instance.TokenName)
		if err != nil {
			p.logger.Warn(fmt.Sprintf("  skipping export, failed to get token with name '%s' from Keepass. Err: %v", instance.TokenName, err))
			continue
		}

		httpClient := httpclient.NewHTTPClient(instance.Hostname, token)
		groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
		dynamicGroupsEndpoint := rest.NewDynamicGroupsEndpoint(httpClient)
		personEndpoint := rest.NewPersonsEndpoint(httpClient)

		for _, group := range instance.Groups {
			p.logger.Info("")
			p.logger.Info(fmt.Sprintf("  processing group '%s'", group.Name))

			persons, err := groupExporter.ExportGroupMembers(
				group.Name,
				groupsEndpoint,
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

			if blocklistsDataProvider.BlockListExists(group) {
				p.logger.Info(fmt.Sprintf("      using blocklist '%s'", group.BlocklistFileName()))
			}

			personData, err := csv.NewPersonData(persons, group, fileDataProvider, blocklistsDataProvider, p.logger)
			if err != nil {
				p.logger.Error(fmt.Sprintf("      failed to extract persons: %v", err))
				continue
			}

			err = os.MkdirAll(filepath.Join(rootDir, instance.Hostname), 0755)
			if err != nil {
				p.logger.Error(fmt.Sprintf("     failed to create directory: %v", err))
				continue
			}

			csvFilePath := filepath.Join(
				rootDir,
				instance.Hostname,
				group.CSVFileName(),
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

func (p instancesProcessor) logTitle(instance config.Instance) {
	boxLength := 70
	title := fmt.Sprintf("Processing instance '%s'", instance.Hostname)
	titleLength := len(title)
	border := strings.Repeat("-", boxLength)

	p.logger.Info("")
	p.logger.Info(fmt.Sprintf("+%s+", border))
	p.logger.Info(fmt.Sprintf("| %s "+strings.Repeat(" ", boxLength-titleLength-2)+"|", title))
	p.logger.Info(fmt.Sprintf("+%s+", border))
}
