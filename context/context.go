package context

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/polydawn/meep"
)

type ResourceSet struct {
	Name   string                 `json:"name"`
	Values map[string]interface{} `json:"values"`
}

type Context struct {
	Name         string                 `json:"context"`
	Global       map[string]interface{} `json:"global"`
	ResourceSets []ResourceSet          `json:"include"`
	BaseDir      string
}

type ContextLoadingError struct {
	meep.AllTraits
	Filename string
}

// Attempt to load and deserialise a Context from the specified file.
func LoadContextFromFile(filename string) (*Context, error) {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, meep.New(
			&ContextLoadingError{Filename: filename},
			meep.Cause(err),
		)
	}

	var c Context

	if strings.HasSuffix(filename, "json") {
		err = json.Unmarshal(file, &c)
	} else if strings.HasSuffix(filename, "yaml") || strings.HasSuffix(filename, "yml") {
		err = yaml.Unmarshal(file, &c)
	} else {
		return nil, meep.New(
			&ContextLoadingError{Filename: filename},
			meep.Cause(fmt.Errorf("File format not supported. Must be JSON or YAML.")),
		)
	}

	if err != nil {
		return nil, meep.New(
			&ContextLoadingError{Filename: filename},
			meep.Cause(err),
		)
	}

	c.BaseDir = path.Dir(filename)

	return &c, nil
}
