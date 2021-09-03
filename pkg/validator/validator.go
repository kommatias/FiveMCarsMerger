package validator

import "strings"

type Validator interface {
	IsValidDataFile(file string) bool
	IsValidStreamFile(file string) bool
}

type validator struct {
}

func New() Validator {
	return &validator{}
}

func (v *validator) IsValidDataFile(file string) bool {
	validDataExtensions := []string{".meta", ".dat"}
	return v.HasAnyFileExtension(file, validDataExtensions)
}

func (v *validator) IsValidStreamFile(file string) bool {
	validStreamExtensions := []string{".yft", ".ytd"}
	return v.HasAnyFileExtension(file, validStreamExtensions)
}

func (v *validator) HasAnyFileExtension(file string, extensions []string) bool {
	for _, extension := range extensions {
		if strings.HasSuffix(file, extension) {
			return true
		}
	}
	return false
}
