package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"
)

// Filenames excluded from templating for the purpose of containing default variable values inside a resource set.
var DefaultFilenames []string = []string{"default.yml", "default.yaml", "default.json"}

// Merges two maps together. Values from the second map override values in the first map.
// The returned map is new if anything was changed.
func Merge(in1 *map[string]interface{}, in2 *map[string]interface{}) *map[string]interface{} {
	if in1 == nil || len(*in1) == 0 {
		return in2
	}

	if in2 == nil || len(*in2) == 0 {
		return in1
	}

	new := make(map[string]interface{})
	for k, v := range *in1 {
		new[k] = v
	}

	for k, v := range *in2 {
		new[k] = v
	}

	return &new
}

// Loads either a YAML or JSON file from the specified path and deserialises it into the provided interface.
func LoadJsonOrYaml(filename string, addr interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if strings.HasSuffix(filename, "json") {
		err = json.Unmarshal(file, addr)
	} else if strings.HasSuffix(filename, "yaml") || strings.HasSuffix(filename, "yml") {
		err = yaml.Unmarshal(file, addr)
	} else {
		return fmt.Errorf("File format not supported. Must be JSON or YAML.")
	}

	return nil
}
