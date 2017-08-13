package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/savaki/aiven/kafka"
	"gopkg.in/urfave/cli.v1"
)

var opts = struct {
	Email    string
	Password string
	OTP      string
	Project  string
	Service  string
	Topic    struct {
		Name           string
		CleanupPolicy  string
		Partitions     int
		Replication    int
		RetentionHours int
	}
}{}

var (
	flagEmail = cli.StringFlag{
		Name:        "email",
		Usage:       "aiven email",
		EnvVar:      "AIVEN_EMAIL",
		Destination: &opts.Email,
	}
	flagPassword = cli.StringFlag{
		Name:        "password",
		Usage:       "aiven password",
		EnvVar:      "AIVEN_PASSWORD",
		Destination: &opts.Password,
	}
	flagOTP = cli.StringFlag{
		Name:        "otp",
		Usage:       "aiven one time password (otp)",
		EnvVar:      "AIVEN_OTP",
		Destination: &opts.OTP,
	}
	flagProject = cli.StringFlag{
		Name:        "project",
		Usage:       "aiven project",
		EnvVar:      "AIVEN_PROJECT",
		Destination: &opts.Project,
	}
	flagService = cli.StringFlag{
		Name:        "service",
		Usage:       "aiven service",
		EnvVar:      "AIVEN_SERVICE",
		Destination: &opts.Service,
	}
	// kafka specific
	//
	flagName = cli.StringFlag{
		Name:        "name",
		Usage:       "name of topic",
		EnvVar:      "TOPIC_NAME",
		Destination: &opts.Topic.Name,
	}
	flagCleanupPolicy = cli.StringFlag{
		Name:        "cleanup-policy",
		Value:       kafka.CleanupPolicyDelete,
		Usage:       "cleanup policy",
		EnvVar:      "TOPIC_CLEANUP_POLICY",
		Destination: &opts.Topic.CleanupPolicy,
	}
	flagPartitions = cli.IntFlag{
		Name:        "partitions",
		Value:       1,
		Usage:       "partitions",
		EnvVar:      "TOPIC_PARTITIONS",
		Destination: &opts.Topic.Partitions,
	}
	flagReplication = cli.IntFlag{
		Name:        "replication",
		Value:       3,
		Usage:       "replication factor",
		EnvVar:      "TOPIC_REPLICATION",
		Destination: &opts.Topic.Replication,
	}
	flagRetentionHours = cli.IntFlag{
		Name:        "retention-hours",
		Value:       36,
		Usage:       "hours to retain content",
		EnvVar:      "TOPIC_REPLICATION_HOURS",
		Destination: &opts.Topic.RetentionHours,
	}
)

func Do(fn func(ctx context.Context) (interface{}, error)) cli.ActionFunc {
	return func(*cli.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		out, err := fn(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if out != nil {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			encoder.Encode(out)
		}

		return nil
	}
}
