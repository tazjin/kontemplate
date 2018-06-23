// Copyright (C) 2016-2017  Vincent Ambo <mail@tazj.in>
//
// Kontemplate is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tazjin/kontemplate/context"
	"github.com/tazjin/kontemplate/templater"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version string = "1.6.0"

// This variable will be initialised by the Go linker during the builder
var gitHash string

var (
	app = kingpin.New("kontemplate", "simple Kubernetes resource templating")

	// Global flags
	includes   = app.Flag("include", "Resource sets to include explicitly").Short('i').Strings()
	excludes   = app.Flag("exclude", "Resource sets to exclude explicitly").Short('e').Strings()
	variables  = app.Flag("var", "Provide variables to templates explicitly").Strings()
	kubectlBin = app.Flag("kubectl", "Path to the kubectl binary (default 'kubectl')").Default("kubectl").String()

	// Commands
	template          = app.Command("template", "Template resource sets and print them")
	templateFile      = template.Arg("file", "Cluster configuration file to use").Required().String()
	templateOutputDir = template.Flag("output", "Output directory in which to save templated files instead of printing them").Short('o').String()

	apply       = app.Command("apply", "Template resources and pass to 'kubectl apply'")
	applyFile   = apply.Arg("file", "Cluster configuration file to use").Required().String()
	applyDryRun = apply.Flag("dry-run", "Print remote operations without executing them").Default("false").Bool()

	replace     = app.Command("replace", "Template resources and pass to 'kubectl replace'")
	replaceFile = replace.Arg("file", "Cluster configuration file to use").Required().String()

	delete     = app.Command("delete", "Template resources and pass to 'kubectl delete'")
	deleteFile = delete.Arg("file", "Cluster configuration file to use").Required().String()

	create     = app.Command("create", "Template resources and pass to 'kubectl create'")
	createFile = create.Arg("file", "Cluster configuration file to use").Required().String()

	versionCmd = app.Command("version", "Show kontemplate version")
)

func main() {
	app.HelpFlag.Short('h')

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case template.FullCommand():
		templateCommand()

	case apply.FullCommand():
		applyCommand()

	case replace.FullCommand():
		replaceCommand()

	case delete.FullCommand():
		deleteCommand()

	case create.FullCommand():
		createCommand()

	case versionCmd.FullCommand():
		versionCommand()
	}
}

func versionCommand() {
	if gitHash == "" {
		fmt.Printf("Kontemplate version %s (git commit unknown)\n", version)
	} else {
		fmt.Printf("Kontemplate version %s (git commit: %s)\n", version, gitHash)
	}
}

func templateCommand() {
	_, resourceSets := loadContextAndResources(templateFile)

	for _, rs := range *resourceSets {
		if len(rs.Resources) == 0 {
			fmt.Fprintf(os.Stderr, "Warning: Resource set '%s' does not exist or contains no valid templates\n", rs.Name)
			continue
		}

		if *templateOutputDir != "" {
			templateIntoDirectory(templateOutputDir, rs)
		} else {
			for _, r := range rs.Resources {
				fmt.Fprintf(os.Stderr, "Rendered file %s/%s:\n", rs.Name, r.Filename)
				fmt.Println(r.Rendered)
			}
		}
	}
}

func templateIntoDirectory(outputDir *string, rs templater.RenderedResourceSet) {
	// Attempt to create the output directory if it does not
	// already exist:
	if err := os.MkdirAll(*templateOutputDir, 0775); err != nil {
		app.Fatalf("Could not create output directory: %v\n", err)
	}

	// Nested resource sets may contain slashes in their names.
	// These are replaced with dashes for the purpose of writing a
	// flat list of output files:
	setName := strings.Replace(rs.Name, "/", "-", -1)

	for _, r := range rs.Resources {
		filename := fmt.Sprintf("%s/%s-%s", *templateOutputDir, setName, r.Filename)
		fmt.Fprintf(os.Stderr, "Writing file %s\n", filename)

		file, err := os.Create(filename)
		if err != nil {
			app.Fatalf("Could not create file %s: %v\n", filename, err)
		}

		_, err = fmt.Fprintf(file, r.Rendered)
		if err != nil {
			app.Fatalf("Error writing file %s: %v\n", filename, err)
		}
	}
}

func applyCommand() {
	ctx, resources := loadContextAndResources(applyFile)

	var kubectlArgs []string

	if *applyDryRun {
		kubectlArgs = []string{"apply", "-f", "-", "--dry-run"}
	} else {
		kubectlArgs = []string{"apply", "-f", "-"}
	}

	if err := runKubectlWithResources(ctx, &kubectlArgs, resources); err != nil {
		failWithKubectlError(err)
	}
}

func replaceCommand() {
	ctx, resources := loadContextAndResources(replaceFile)
	args := []string{"replace", "--save-config=true", "-f", "-"}

	if err := runKubectlWithResources(ctx, &args, resources); err != nil {
		failWithKubectlError(err)
	}
}

func deleteCommand() {
	ctx, resources := loadContextAndResources(deleteFile)
	args := []string{"delete", "-f", "-"}

	if err := runKubectlWithResources(ctx, &args, resources); err != nil {
		failWithKubectlError(err)
	}
}

func createCommand() {
	ctx, resources := loadContextAndResources(createFile)
	args := []string{"create", "--save-config=true", "-f", "-"}

	if err := runKubectlWithResources(ctx, &args, resources); err != nil {
		failWithKubectlError(err)
	}
}

func loadContextAndResources(file *string) (*context.Context, *[]templater.RenderedResourceSet) {
	ctx, err := context.LoadContext(*file, variables)
	if err != nil {
		app.Fatalf("Error loading context: %v\n", err)
	}

	resources, err := templater.LoadAndApplyTemplates(includes, excludes, ctx)
	if err != nil {
		app.Fatalf("Error templating resource sets: %v\n", err)
	}

	return ctx, &resources
}

func runKubectlWithResources(c *context.Context, kubectlArgs *[]string, resourceSets *[]templater.RenderedResourceSet) error {
	args := append(*kubectlArgs, fmt.Sprintf("--context=%s", c.Name))

	for _, rs := range *resourceSets {
		if len(rs.Resources) == 0 {
			fmt.Fprintf(os.Stderr, "Warning: Resource set '%s' contains no valid templates\n", rs.Name)
			continue
		}

		kubectl := exec.Command(*kubectlBin, args...)

		stdin, err := kubectl.StdinPipe()
		if err != nil {
			return fmt.Errorf("kubectl error: %v", err)
		}

		kubectl.Stdout = os.Stdout
		kubectl.Stderr = os.Stderr

		if err = kubectl.Start(); err != nil {
			return fmt.Errorf("kubectl error: %v", err)
		}

		for _, r := range rs.Resources {
			fmt.Printf("Passing file %s/%s to kubectl\n", rs.Name, r.Filename)
			fmt.Fprintln(stdin, r.Rendered)
		}
		stdin.Close()

		if err = kubectl.Wait(); err != nil {
			return err
		}
	}

	return nil
}

func failWithKubectlError(err error) {
	app.Fatalf("Kubectl error: %v\n", err)
}
