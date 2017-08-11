package lib

import (
	"context"

	"github.com/savaki/aiven"
	"github.com/savaki/aiven/kafka"
	"gopkg.in/urfave/cli.v1"
)

var Kafka = cli.Command{
	Name:  "kafka",
	Usage: "kafka related commands",
	Subcommands: cli.Commands{
		{
			Name:  "topics",
			Usage: "list kafka topics",
			Flags: []cli.Flag{
				flagEmail,
				flagPassword,
				flagOTP,
				flagProject,
				flagService,
			},
			Action: Do(kafkaTopics),
		},
	},
}

func kafkaTopics(ctx context.Context) (interface{}, error) {
	client, err := aiven.NewOTP(opts.Email, opts.Password, opts.OTP)
	if err != nil {
		return nil, err
	}

	return kafka.New(client).Topics(ctx, opts.Project, opts.Service)
}
