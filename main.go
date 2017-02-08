package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/polydawn/meep"
	"github.com/tazjin/kontemplate/context"
	"github.com/tazjin/kontemplate/templater"
	"github.com/urfave/cli"
)

type KubeCtlError struct {
	meep.AllTraits
}

func main() {
	app := cli.NewApp()

	app.Name = "kontemplate"
	app.Usage = "simple Kubernetes resource templating"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		templateCommand(),
		applyCommand(),
		replaceCommand(),
		deleteCommand(),
	}

	app.Run(os.Args)
}

func templateCommand() cli.Command {
	return cli.Command{
		Name:  "template",
		Usage: "Interpolate and print templates",
		Flags: commonFlags(),
		Action: func(c *cli.Context) error {
			include := c.StringSlice("include")
			exclude := c.StringSlice("exclude")
			ctx, err := loadContext(c)
			resources, err := templater.LoadAndPrepareTemplates(&include, &exclude, ctx)

			if err != nil {
				return err
			}

			for _, r := range resources {
				fmt.Println(r)
			}

			return nil
		},
	}
}

func applyCommand() cli.Command {
	dryRun := false

	return cli.Command{
		Name:  "apply",
		Usage: "Interpolate templates and run 'kubectl apply'",
		Flags: append(commonFlags(), cli.BoolFlag{
			Name:        "dry-run",
			Usage:       "Only print objects that would be sent, without sending them",
			Destination: &dryRun,
		}),
		Action: func(c *cli.Context) error {
			include := c.StringSlice("include")
			exclude := c.StringSlice("exclude")
			ctx, err := loadContext(c)
			resources, err := templater.LoadAndPrepareTemplates(&include, &exclude, ctx)

			if err != nil {
				return err
			}

			var args []string
			if dryRun {
				args = []string{"apply", "-f", "-", "--dry-run"}
			} else {
				args = []string{"apply", "-f", "-"}
			}

			return runKubectlWithResources(ctx, &args, &resources)
		},
	}
}

func replaceCommand() cli.Command {
	return cli.Command{
		Name:  "replace",
		Usage: "Interpolate templates and run 'kubectl replace'",
		Flags: commonFlags(),
		Action: func(c *cli.Context) error {
			include := c.StringSlice("include")
			exclude := c.StringSlice("exclude")
			ctx, err := loadContext(c)
			resources, err := templater.LoadAndPrepareTemplates(&include, &exclude, ctx)

			if err != nil {
				return err
			}

			args := []string{"replace", "--save-config=true", "-f", "-"}
			return runKubectlWithResources(ctx, &args, &resources)
		},
	}
}

func deleteCommand() cli.Command {
	return cli.Command{
		Name:  "delete",
		Usage: "Interpolate templates and run 'kubectl delete'",
		Flags: commonFlags(),
		Action: func(c *cli.Context) error {
			include := c.StringSlice("include")
			exclude := c.StringSlice("exclude")
			ctx, err := loadContext(c)
			resources, err := templater.LoadAndPrepareTemplates(&include, &exclude, ctx)

			if err != nil {
				return err
			}

			args := []string{"delete", "-f", "-"}
			return runKubectlWithResources(ctx, &args, &resources)
		},
	}
}

func runKubectlWithResources(c *context.Context, kubectlArgs *[]string, resources *[]string) error {
	args := append(*kubectlArgs, fmt.Sprintf("--context=%s", c.Name))

	kubectl := exec.Command("kubectl", args...)

	stdin, err := kubectl.StdinPipe()
	if err != nil {
		return meep.New(&KubeCtlError{}, meep.Cause(err))
	}

	kubectl.Stdout = os.Stdout
	kubectl.Stderr = os.Stderr

	if err = kubectl.Start(); err != nil {
		return meep.New(&KubeCtlError{}, meep.Cause(err))
	}

	for _, r := range *resources {
		fmt.Fprintln(stdin, r)
	}
	stdin.Close()

	kubectl.Wait()

	return nil
}

func commonFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Cluster configuration file to use",
		},
		cli.StringSliceFlag{
			Name:  "include, i",
			Usage: "Limit templating to explicitly included resource sets",
		},
		cli.StringSliceFlag{
			Name:  "exclude, e",
			Usage: "Exclude certain resource sets from templating",
		},
	}
}

func loadContext(c *cli.Context) (*context.Context, error) {
	f := c.String("file")

	if f == "" {
		return nil, meep.New(
			&meep.ErrInvalidParam{
				Param:  "file",
				Reason: "Cluster config file must be specified",
			},
		)
	}

	ctx, err := context.LoadContextFromFile(f)

	if err != nil {
		return nil, err
	}

	return ctx, nil
}
