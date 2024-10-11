package app

import (
    "ctRestClient/config"
    "ctRestClient/httpclient"
    "ctRestClient/rest"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "time"
)

type InstancesProcessor interface {
    Process(groupExporter GroupExporter, csvWriter CSVWriter) error
}

type instancesProcessor struct {
    config          config.Config
    outputDirectory string
    logger		  Logger
}

func NewInstancesProcessor(
    config config.Config,
    outputDirectory string,
    logger Logger,
) InstancesProcessor {
    return instancesProcessor{
        config:          config,
        outputDirectory: outputDirectory,
        logger:          logger,
    }
}

func (p instancesProcessor) Process(groupExporter GroupExporter, csvWriter CSVWriter) error {
    // define the root dir location via parameter that is by default beside the executable
    rootDir := filepath.Join(p.outputDirectory, "export", time.Now().Format("2006.01.02_15-04-05"))

    for _, instance := range p.config.Instances {
        p.logger.Info(fmt.Sprintf("processing instance '%s'", instance.Hostname))
        
        token := os.Getenv(instance.TokenName)
        if token == "" {
            p.logger.Warn(fmt.Sprintf("  skipping export, a token with name '%s' is not set in the environment", instance.TokenName))
            continue
        }

        httpClient := httpclient.NewHTTPClient(instance.Hostname, token)
        dynamicGroupsEndpoint := rest.NewDynamicGroupsEndpoint(httpClient)
        groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
        personEndpoint := rest.NewPersonsEndpoint(httpClient)

        groupName2IDMap, err := groupExporter.GetGroupNames2IDMapping(dynamicGroupsEndpoint, groupsEndpoint)
        if err != nil {
            return fmt.Errorf("failed to resolve groupnames to ids, %w", err)
        }

        for _, group := range instance.Groups {
            p.logger.Info(fmt.Sprintf("  processing group '%s'", group.Name))

            err := os.MkdirAll(filepath.Join(rootDir, instance.Hostname), 0755)
            if err != nil {
                return fmt.Errorf("failed to create group directory: %v", err)
            }

            groupID, ok := groupName2IDMap[group.Name]
            if !ok {
                p.logger.Error("    could not find group to id mapping")
                continue
            }
  
            persons, err := groupExporter.ExportPersonData(
                groupID,
                groupsEndpoint,
                personEndpoint,
            )
            if err != nil {
                return fmt.Errorf("failed to get person informations: %v", err)
            }

            if len(persons) == 0 {
                p.logger.Info("    the group is empty")
                continue
            } else {
                p.logger.Info(fmt.Sprintf("    the group has %d persons", len(persons)))
            }

            csvHeader := group.Fields
            csvRecords := make([][]string, 0)

            for _, person := range persons {
                var data map[string]interface{}
                err := json.Unmarshal([]byte(person), &data)
                if err != nil {
                    return fmt.Errorf("failed to read person information raw json: %v", err)
                }

                record := make([]string, len(group.Fields))

                for i, field := range group.Fields {
                    if value, ok := data[field].(string); ok {
                        // get string values
                        record[i] = value
                    } else if value, ok := data[field].(float64); ok {
                        // get int values
                        record[i] = fmt.Sprintf("%d", int(value))
                    } else {
                        p.logger.Warn(fmt.Sprintf("    Field %s is not a string or int, or not found", field))
                        record[i] = ""
                    }
                }
                csvRecords = append(csvRecords, record)
            }

            csvFilePath := filepath.Join(
                rootDir,
                instance.Hostname,
                group.SanitizedGroupName()+".csv",
            )

            err = csvWriter.Write(csvFilePath, csvHeader, csvRecords)
            if err != nil {
                return fmt.Errorf("failed to write csv file: %v", err)
            }
        }
    }

    return nil
}
