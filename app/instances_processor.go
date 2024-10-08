package app

import (
    "ctRestClient/config"
    "ctRestClient/httpclient"
    "ctRestClient/rest"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"
)

type InstancesProcessor interface {
    Process(groupExporter GroupExporter, csvWriter CSVWriter) error
}

type instancesProcessor struct {
    config config.Config
    token  string
    outputDirectory string
}

func NewInstancesProcessor(config config.Config, token string, outputDirectory string) InstancesProcessor {
    return instancesProcessor{
        config: config,
        token:  token,
        outputDirectory: outputDirectory,
    }
}

func (p instancesProcessor) Process(groupExporter GroupExporter, csvWriter CSVWriter) error {
    // define the root dir location via parameter that is by default beside the executable
    rootDir := filepath.Join(p.outputDirectory, "export", time.Now().Format("2006.01.02_15-04-05"))

    for _, instance := range p.config.Instances {
        for _, group := range instance.Groups {
            log.Printf("[INFO] processing '%s' group '%s'", instance.HostName, group.Name)

            err := os.MkdirAll(filepath.Join(rootDir, instance.HostName), 0755)
            if err != nil {
                return fmt.Errorf("failed to create group directory: %v", err)
            }

            httpClient := httpclient.NewHTTPClient(instance.HostName, p.token)
            dynamicGroupsEndpoint := rest.NewDynamicGroupsEndpoint(httpClient)
            groupsEndpoint := rest.NewGroupsEndpoint(httpClient)
            personEndpoint := rest.NewPersonsEndpoint(httpClient)

            persons, err := groupExporter.ExportPersonData(
                group.Name,
                dynamicGroupsEndpoint,
                groupsEndpoint,
                personEndpoint,
            )
            if err != nil {
                return fmt.Errorf("failed to get person informations: %v", err)
            }

            if len(persons) == 0 {
                log.Printf("[INFO]   the group is empty")
                continue
            } else {
                log.Printf("[INFO]   the group has %d persons", len(persons))
            }

            csvHeader := group.Fields
            csvRecords := make([][]string, 0)

            for _, person := range persons {
                // log.Printf("[INFO]   processing person: %s", person)
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
                        log.Printf("Field %s is not a string or int, or not found", field)
                        record[i] = ""
                    }
                }
                csvRecords = append(csvRecords, record)
            }

            csvFilePath := filepath.Join(
                rootDir,
                instance.HostName,
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
