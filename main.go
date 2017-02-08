package main

import (
	"fmt"
	"os"

	"github.com/polydawn/meep"
	"github.com/tazjin/kontemplate/context"
	"github.com/tazjin/kontemplate/templater"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "kontemplate"
	app.Usage = "simple Kubernetes resource templating"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		ApplyCommand(),
	}

	app.Run(os.Args)
}

func ApplyCommand() cli.Command {
	return cli.Command{
		Name:  "run",
		Usage: "Interpolate and print templates",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "Cluster configuration file to use",
			},
			cli.StringSliceFlag{
				Name:  "limit, l",
				Usage: "Limit templating to certain resource sets",
			},
		},
		Action: func(c *cli.Context) error {
			limit := c.StringSlice("limit")
			f := c.String("file")

			if f == "" {
				return meep.New(
					&meep.ErrInvalidParam{
						Param:  "file",
						Reason: "Cluster config file must be specified",
					},
				)
			}

			ctx, err := context.LoadContextFromFile(f)

			if err != nil {
				return err
			}

			resources, err := templater.LoadAndPrepareTemplates(&limit, ctx)

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
