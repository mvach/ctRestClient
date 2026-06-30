package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ExportGroupConfig struct {
	Instances []ExportGroupInstance `yaml:"instances"`
}

func (g ExportGroupConfig) LoadConfig(filePath string) (Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config ExportGroupConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to load invalid config file, %w", err)
	}

	err = g.validate(config)
	if err != nil {
		return nil, fmt.Errorf("failed to validate the config file, %w", err)
	}

	return config, nil
}

func (g ExportGroupConfig) validate(config ExportGroupConfig) error {
	if len(config.Instances) == 0 {
		return errors.New("property instances is not set")
	}

	for _, instance := range config.Instances {
		err := instance.Validate()
		if err != nil {
			return fmt.Errorf("failed to validate instance %s, %w", instance.Hostname, err)
		}
	}
	return nil
}
