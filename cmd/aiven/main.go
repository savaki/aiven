package main

import (
	"os"

	"github.com/savaki/aiven/cmd/aiven/lib"
	"gopkg.in/urfave/cli.v1"
)

const (
	Version = "SNAPSHOT"
)

type Options struct {
}

var opts Options

func main() {
	app := cli.NewApp()
	app.Usage = "console interface to aiven"
	app.Version = Version
	app.Commands = cli.Commands{
		lib.Kafka,
	}
	app.Run(os.Args)
}

func Run(_ *cli.Context) error {
	return nil
}
