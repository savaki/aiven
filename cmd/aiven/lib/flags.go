package lib

import (
	"fmt"
	"os"

	"encoding/json"

	"context"
	"time"

	"gopkg.in/urfave/cli.v1"
)

var opts = struct {
	Email    string
	Password string
	OTP      string
	Project  string
	Service  string
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

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.Encode(out)

		return nil
	}
}
