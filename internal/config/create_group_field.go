package config

import (
	"fmt"
)

type CreateGroupField struct {
	FieldName   string    `yaml:"fieldname"`
	FieldType   FieldType `yaml:"fieldType"`
	Description string    `yaml:"description"`
	DataType    DataType  `yaml:"dataType"`
	Mandatory   bool      `yaml:"mandatory"`
}

func (e *CreateGroupField) IsField() {}

type FieldType string

const (
	PersonField DataType = "person"
	GroupField  DataType = "group"
)

type DataType string

const (
	Text    DataType = "text"
	Boolean DataType = "boolean"
	Date    DataType = "date"
)

func (c *CreateGroupField) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Define a temporary type to capture the raw YAML structure
	type rawCreateField struct {
		FieldName   string    `yaml:"fieldname"`
		FieldType   FieldType `yaml:"fieldType"`
		Description string    `yaml:"description"`
		DataType    DataType  `yaml:"dataType"`
		Mandatory   bool      `yaml:"mandatory"`
	}

	raw := rawCreateField{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	c.FieldName = raw.FieldName
	c.FieldType = raw.FieldType
	c.Description = raw.Description
	c.DataType = raw.DataType
	c.Mandatory = raw.Mandatory

	// Validate required fields for CreateGroupField
	if c.FieldName == "" || c.FieldType == "" || c.DataType == "" {
		return fmt.Errorf("'fieldname', 'fieldType', and 'dataType' must be set for CreateGroupField")
	}

	return nil
}