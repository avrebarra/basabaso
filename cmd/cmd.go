package cmd

import (
	"os"

	"github.com/avrebarra/valeed"
	"github.com/urfave/cli"
)

type Config struct{}

type Command struct {
	config  Config
	cmdRoot *cli.App
}

func New(cfg Config) (*Command, error) {
	if err := valeed.Validate(cfg); err != nil {
		return nil, err
	}

	e := &Command{config: cfg}
	e.init()

	return e, nil
}

func (c *Command) init() {
	c.cmdRoot = &cli.App{
		Name:    "basabaso",
		Version: "v1",
		Usage:   "golang service",
		Action: func(c *cli.Context) error {
			return ExecDefault{}.Run()
		},
		Commands: []cli.Command{
			{
				Name:  "server",
				Usage: "start application server",
				Action: func(c *cli.Context) error {
					return ExecDefault{}.Run()
				},
			},
		},
	}

}

func (c *Command) Execute() error {
	return c.cmdRoot.Run(os.Args)
}
