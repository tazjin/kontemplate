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
	}

	app.Run(os.Args)
}

func templateCommand() cli.Command {
	return cli.Command{
		Name:  "template",
		Usage: "Interpolate and print templates",
		Flags: commonFlags(),
		Action: func(c *cli.Context) error {
			resources, err := templateResources(c)

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
			resources, err := templateResources(c)

			if err != nil {
				return err
			}

			var kubectl *exec.Cmd
			if dryRun {
				kubectl = exec.Command("kubectl", "apply", "-f", "-", "--dry-run")
			} else {
				kubectl = exec.Command("kubectl", "apply", "-f", "-")
			}

			stdin, err := kubectl.StdinPipe()
			if err != nil {
				return meep.New(&KubeCtlError{}, meep.Cause(err))
			}

			kubectl.Stdout = os.Stdout
			kubectl.Stderr = os.Stderr

			if err = kubectl.Start(); err != nil {
				return meep.New(&KubeCtlError{}, meep.Cause(err))
			}

			for _, r := range resources {
				fmt.Fprintln(stdin, r)
			}

			stdin.Close()

			kubectl.Wait()

			return nil
		},
	}
}

func commonFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Cluster configuration file to use",
		},
		cli.StringSliceFlag{
			Name:  "limit, l",
			Usage: "Limit templating to certain resource sets",
		},
	}
}

func templateResources(c *cli.Context) ([]string, error) {
	limit := c.StringSlice("limit")
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

	return templater.LoadAndPrepareTemplates(&limit, ctx)
}
