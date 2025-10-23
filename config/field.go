package config

import "fmt"

// A Field can be either a simple string or a structured object
// with field and column names.
type Field struct {
	FieldName *string
	Object    *FieldInformation
}

func (f *Field) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try unmarshaling as a simple string
	var s string
	if err := unmarshal(&s); err == nil {
		f.FieldName = &s
		return nil
	}

	// Try unmarshaling as a structured object
	var obj FieldInformation
	if err := unmarshal(&obj); err == nil {
		// Validate required fields
		if obj.FieldName == "" || obj.ColumnName == "" {
			return fmt.Errorf("both 'fieldname' and 'columnname' must be set")
		}
		f.Object = &obj
		return nil
	}

	return fmt.Errorf("field must be a string or an object with 'fieldname' and 'columnname'")
}

func (f *Field) GetFieldName() string {
	if f.Object != nil {
		return f.Object.FieldName
	}
	if f.FieldName != nil {
		return *f.FieldName
	}
	return ""
}

func (f *Field) GetColumnName() string {
	if f.Object != nil {
		return f.Object.ColumnName
	}
	if f.FieldName != nil {
		return *f.FieldName
	}
	return ""
}

func (f *Field) IsMappedData() bool {
	if f.FieldName == nil && f.Object != nil {
		return true
	}
	return false
}
