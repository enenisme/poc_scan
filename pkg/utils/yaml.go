// pkg/utils/yaml.go
package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/enenisme/poc_scan/pkg/model"
	"gopkg.in/yaml.v2"
)

func LoadPOCFromFile(path string) (*model.POC, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read POC file: %v", err)
	}

	var poc model.POC
	err = yaml.Unmarshal(data, &poc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %v", err)
	}

	// Validate required fields
	if poc.Name == "" {
		return nil, fmt.Errorf("POC name is required")
	}
	if poc.Transport != "http" {
		return nil, fmt.Errorf("only HTTP transport is supported")
	}

	return &poc, nil
}
