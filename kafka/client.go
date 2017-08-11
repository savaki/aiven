package kafka

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/savaki/aiven"
)

// Api provides an api into aiven kafka
type Api struct {
	client *aiven.Client
}

// New accepts a valid aiven client and returns access to the Kafka api
func New(client *aiven.Client) *Api {
	return &Api{
		client: client,
	}
}

// Topic represents the Kafka topics
type Topic struct {
	CleanupPolicy  string `json:"cleanup_policy"`
	Partitions     int    `json:"partitions"`
	Replication    int    `json:"replication"`
	RetentionHours int    `json:"retention_hours"`
	State          string `json:"state"`
	TopicName      string `json:"topic_name"`
}

// Topics returns the list of topics for the service
func (a *Api) Topics(ctx context.Context, project, service string) ([]Topic, error) {
	u := fmt.Sprintf("https://console.aiven.io/v1beta/project/%v/service/%v", project, service)
	out := struct {
		Service struct {
			Topics []Topic
		}
	}{}
	if err := a.client.Get(ctx, u, &out); err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve topics for project:service, %v:%v", project, service)
	}

	return out.Service.Topics, nil
}
