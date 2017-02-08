package main

import (
	"fmt"
	"os"

	"github.com/tazjin/kontemplate/context"
	"github.com/tazjin/kontemplate/templater"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: kontemplate <cluster-config>")
		os.Exit(1)
	}

	c, err := context.LoadContextFromFile(os.Args[1])

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Applying cluster %s\n", c.Name)

	for _, rs := range c.ResourceSets {
		fmt.Fprintf(os.Stderr, "Applying resource %s with values %v\n", rs.Name, rs.Values)
		resources, err := templater.LoadAndPrepareTemplates(c)

		if err != nil {
			fmt.Println(err)
		}

		for _, r := range resources {
			fmt.Print(r)
		}
	}
}
