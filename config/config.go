package config

import (
    "errors"
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

type Group struct {
    Name    string   `yaml:"name"`
    MergeBy string   `yaml:"merge_by,omitempty"`
    Fields  []string `yaml:"fields"`
}

type Instance struct {
    Instance string  `yaml:"instance"`
    Groups  []Group `yaml:"groups"`
}

type Config struct {
    Instances []Instance `yaml:"instances"`
}

func LoadConfig(filePath string) (*Config, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, fmt.Errorf("failed to load invalid config file, %w", err)
    }

    err = config.validate()
    if err != nil {
        return nil, fmt.Errorf("failed to validate the config file, %w", err)
    }

    return &config, nil
}

func (c *Config) validate() error {
    if len(c.Instances) == 0 {
        return errors.New("property instances is not set")
    }
    for _, instance := range c.Instances {
        if instance.Instance == "" {
            return errors.New("property instance is not set")
        }

        if len(instance.Groups) == 0 {
            return errors.New("property groups is not set")
        }
        for _, group := range instance.Groups {
            if group.Name == "" {
                return errors.New("property name is not set")
            }
            if len(group.Fields) == 0 {
                return errors.New("property fields is not set")
            }
        }
    }
    return nil
}