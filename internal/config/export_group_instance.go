package config

import (
	"errors"
)

type ExportGroupInstance struct {
	Hostname  string             `yaml:"hostname"`
	TokenName string             `yaml:"token_name"`
	Groups    []ExportGroupGroup `yaml:"groups"`
}

func (c ExportGroupInstance) Validate() error {
	if c.Hostname == "" {
		return errors.New("property hostname is not set")
	}
	if c.TokenName == "" {
		return errors.New("property token_name is not set")
	}
	if len(c.Groups) == 0 {
		return errors.New("property groups is not set")
	}
	for _, group := range c.Groups {
		err := group.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
