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
			Name:  "list-topics",
			Usage: "list kafka topics",
			Flags: []cli.Flag{
				flagEmail,
				flagPassword,
				flagOTP,
				flagProject,
				flagService,
			},
			Action: Do(listTopics),
		},
		{
			Name:  "create-topic",
			Usage: "create kafka topic",
			Flags: []cli.Flag{
				flagEmail,
				flagPassword,
				flagOTP,
				flagProject,
				flagService,
				flagName,
				flagCleanupPolicy,
				flagPartitions,
				flagReplication,
				flagRetentionHours,
			},
			Action: Do(createTopic),
		},
		{
			Name:  "delete-topic",
			Usage: "delete kafka topic",
			Flags: []cli.Flag{
				flagEmail,
				flagPassword,
				flagOTP,
				flagProject,
				flagService,
				flagName,
			},
			Action: Do(deleteTopic),
		},
	},
}

func listTopics(ctx context.Context) (interface{}, error) {
	client, err := aiven.NewOTP(opts.Email, opts.Password, opts.OTP)
	if err != nil {
		return nil, err
	}

	return kafka.New(client).ListTopics(ctx, kafka.ListTopicsIn{
		Project: opts.Project,
		Service: opts.Service,
	})
}

func createTopic(ctx context.Context) (interface{}, error) {
	client, err := aiven.NewOTP(opts.Email, opts.Password, opts.OTP)
	if err != nil {
		return nil, err
	}

	return nil, kafka.New(client).CreateTopic(ctx, kafka.CreateTopicIn{
		Project:        opts.Project,
		Service:        opts.Service,
		CleanupPolicy:  opts.Topic.CleanupPolicy,
		Partitions:     opts.Topic.Partitions,
		Replication:    opts.Topic.Replication,
		RetentionHours: opts.Topic.RetentionHours,
		TopicName:      opts.Topic.Name,
	})
}

func deleteTopic(ctx context.Context) (interface{}, error) {
	client, err := aiven.NewOTP(opts.Email, opts.Password, opts.OTP)
	if err != nil {
		return nil, err
	}

	return nil, kafka.New(client).DeleteTopic(ctx, kafka.DeleteTopicIn{
		Project:   opts.Project,
		Service:   opts.Service,
		TopicName: opts.Topic.Name,
	})
}
