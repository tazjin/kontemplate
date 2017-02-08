package templater

import (
	"fmt"
	"io/ioutil"
	"strings"
	"os"
	"path"
	"text/template"
	"bytes"

	"github.com/tazjin/kontemplate/context"
	"github.com/polydawn/meep"
)

// Error that is caused by non-existent template files being specified
type TemplateNotFoundError struct {
	meep.AllTraits
	Name string
}

// Error that is caused during templating, e.g. required value being absent or invalid template format
type TemplatingError struct {
	meep.AllTraits
}

func LoadAndPrepareTemplates(c *context.Context) ([]string, error) {
	output := make([]string, 0)

	for _, rs := range c.ResourceSets {
		fmt.Fprintf(os.Stderr,"Loading resources for %s\n", rs.Name)

		rp := path.Join(c.BaseDir, rs.Name)
		files, err := ioutil.ReadDir(rp)

		if err != nil {
			return nil, meep.New(
				&TemplateNotFoundError{Name: rs.Name},
				meep.Cause(err),
			)
		}


		for _, file := range files {
			if !file.IsDir() && isResourceFile(file) {
				p := path.Join(rp, file.Name())
				o, err := templateFile(c, &rs, p)

				if err != nil {
					return nil, err
				}

				output = append(output, o)
			}
		}
	}

	return output, nil
}

func templateFile(c *context.Context, rs *context.ResourceSet, filename string) (string, error) {
	tpl, err := template.ParseFiles(filename)

	if err != nil {
		return "", meep.New(
			&TemplateNotFoundError{Name: filename},
			meep.Cause(err),
		)
	}

	var b bytes.Buffer

	// Merge global and resourceset-specific values (don't override from global)
	for k, v := range c.Global {
		if _, ok := rs.Values[k]; !ok {
			rs.Values[k] = v
		}
	}

	err = tpl.Execute(&b, rs.Values)

	if err != nil {
		return "", meep.New(
			&TemplatingError{},
			meep.Cause(err),
		)
	}

	return b.String(), nil
}

// Checks whether a file is a resource file (i.e. is YAML or JSON)
func isResourceFile(f os.FileInfo) bool {
	return strings.HasSuffix(f.Name(), "yaml") ||
		strings.HasSuffix(f.Name(), "yml") ||
		strings.HasSuffix(f.Name(), "json")
}
