package templater

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/polydawn/meep"
	"github.com/tazjin/kontemplate/context"
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

func LoadAndPrepareTemplates(limit *[]string, c *context.Context) (output []string, err error) {
	for _, rs := range c.ResourceSets {
		if resourceSetIncluded(limit, &rs.Name) {
			err = processResourceSet(c, &rs, &output)

			if err != nil {
				return
			}
		}
	}

	return
}

func resourceSetIncluded(limit *[]string, resourceSetName *string) bool {
	if len(*limit) == 0 {
		return true
	}

	for _, name := range *limit {
		if name == *resourceSetName {
			return true
		}
	}

	return false
}

func processResourceSet(c *context.Context, rs *context.ResourceSet, output *[]string) error {
	fmt.Fprintf(os.Stderr, "Loading resources for %s\n", rs.Name)

	rp := path.Join(c.BaseDir, rs.Name)
	files, err := ioutil.ReadDir(rp)

	err = processFiles(c, rs, rp, files, output)

	if err != nil {
		return meep.New(
			&TemplateNotFoundError{Name: rs.Name},
			meep.Cause(err),
		)
	}

	return nil
}

func processFiles(c *context.Context, rs *context.ResourceSet, rp string, files []os.FileInfo, output *[]string) error {
	for _, file := range files {
		if !file.IsDir() && isResourceFile(file) {
			p := path.Join(rp, file.Name())
			o, err := templateFile(c, rs, p)

			if err != nil {
				return err
			}

			*output = append(*output, o)
		}
	}

	return nil
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
