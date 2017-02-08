package templater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
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

func LoadAndPrepareTemplates(include *[]string, exclude *[]string, c *context.Context) (output []string, err error) {
	limitedResourceSets := applyLimits(&c.ResourceSets, include, exclude)

	for _, rs := range *limitedResourceSets {
		err = processResourceSet(c, &rs, &output)

		if err != nil {
			return
		}
	}

	return
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
	tpl, err := template.New(path.Base(filename)).Funcs(templateFuncs()).ParseFiles(filename)

	if err != nil {
		return "", meep.New(
			&TemplateNotFoundError{Name: filename},
			meep.Cause(err),
		)
	}

	var b bytes.Buffer

	// Guard against empty map before merging keys
	if rs.Values == nil {
		rs.Values = make(map[string]interface{}, 0)
	}

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

// Applies the limits of explicitly included or excluded resources and returns the updated resource set.
// Exclude takes priority over include
func applyLimits(rs *[]context.ResourceSet, include *[]string, exclude *[]string) *[]context.ResourceSet {
	if len(*include) == 0 && len(*exclude) == 0 {
		return rs
	}

	// Exclude excluded resource sets
	excluded := make([]context.ResourceSet, 0)
	for _, r := range *rs {
		if !matchesResourceSet(exclude, &r) {
			excluded = append(excluded, r)
		}
	}

	// Include included resource sets
	if len(*include) == 0 {
		return &excluded
	}
	included := make([]context.ResourceSet, 0)
	for _, r := range excluded {
		if matchesResourceSet(include, &r) {
			included = append(included, r)
		}
	}

	return &included
}

// Check whether an include/exclude string slice matches a resource set
func matchesResourceSet(s *[]string, rs *context.ResourceSet) bool {
	for _, r := range *s {
		if r == rs.Name || r == *rs.Parent {
			return true
		}
	}

	return false
}

func templateFuncs() template.FuncMap {
	m := sprig.TxtFuncMap()
	m["json"] = func(data interface{}) string {
		b, _ := json.Marshal(data)
		return string(b)
	}

	return m
}

// Checks whether a file is a resource file (i.e. is YAML or JSON)
func isResourceFile(f os.FileInfo) bool {
	return strings.HasSuffix(f.Name(), "yaml") ||
		strings.HasSuffix(f.Name(), "yml") ||
		strings.HasSuffix(f.Name(), "json")
}
