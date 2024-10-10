package config

import (
    "errors"
    "fmt"
    "os"
    "regexp"
    "strings"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Instances []Instance `yaml:"instances"`
}

type Instance struct {
    Hostname  string  `yaml:"hostname"`
    TokenName string  `yaml:"token_name"`
    Groups    []Group `yaml:"groups"`
}

type Group struct {
    Name    string   `yaml:"name"`
    MergeBy string   `yaml:"merge_by,omitempty"`
    Fields  []string `yaml:"fields"`
}

func (g Group) SanitizedGroupName() string {
    fileName := g.Name
    fileName = strings.ReplaceAll(fileName, " ", "_")
    fileName = strings.ReplaceAll(fileName, ",", ".")
    fileName = strings.ReplaceAll(fileName, "ä", "ae")
    fileName = strings.ReplaceAll(fileName, "ö", "oe")
    fileName = strings.ReplaceAll(fileName, "ü", "ue")
    fileName = strings.ReplaceAll(fileName, "Ä", "Ae")
    fileName = strings.ReplaceAll(fileName, "Ö", "Oe")
    fileName = strings.ReplaceAll(fileName, "Ü", "Ue")

    re := regexp.MustCompile(`[^\w\-.]`)
    fileName = re.ReplaceAllString(fileName, "")
    fileName = strings.ReplaceAll(fileName, "__", "_")

    return fileName
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

func (c Config) validate() error {
    if len(c.Instances) == 0 {
        return errors.New("property instances is not set")
    }
    for _, instance := range c.Instances {
        if instance.Hostname == "" {
            return errors.New("property hostname is not set")
        }
        if instance.TokenName == "" {
            return errors.New("property token_name is not set")
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
