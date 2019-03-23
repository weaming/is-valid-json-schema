package main

import (
	"strings"

	gojs "github.com/xeipuuv/gojsonschema"
)

type SchemaError struct {
	errors    []string
	IsDataErr bool
}

func (this SchemaError) Error() string {
	return strings.Join(this.errors, "\n")
}

func ValidateSchema(schemaLoader, documentLoader gojs.JSONLoader) (bool, *SchemaError) {
	result, err := gojs.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, &SchemaError{[]string{err.Error()}, false}
	}

	if result.Valid() {
		return true, nil
	} else {
		errs := []string{}
		for _, e := range result.Errors() {
			errs = append(errs, e.String())
		}
		return false, &SchemaError{errs, true}
	}
}
