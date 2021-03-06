package kafka

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/savaki/aiven"
)

// Api provides an api into aiven kafka
type Api struct {
	client *aiven.Client
}

// New accepts a valid aiven client and returns access to the Kafka api
// Deprecated use (*aiven.Client).Kafka() instead
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

type ListTopicsIn struct {
	Project string
	Service string
}

// ListTopics returns the list of all topics
func (a *Api) ListTopics(ctx context.Context, in ListTopicsIn) ([]Topic, error) {
	u := fmt.Sprintf("https://console.aiven.io/v1beta/project/%v/service/%v", in.Project, in.Service)
	out := struct {
		Service struct {
			Topics []Topic
		}
	}{}
	if err := a.client.Get(ctx, u, &out); err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve topics for project:service, %v:%v", in.Project, in.Service)
	}

	return out.Service.Topics, nil
}

type CreateTopicIn struct {
	Project        string `json:"-"`
	Service        string `json:"-"`
	CleanupPolicy  string `json:"cleanup_policy"`
	Partitions     int    `json:"partitions"`
	Replication    int    `json:"replication"`
	RetentionHours int    `json:"retention_hours"`
	TopicName      string `json:"topic_name"`
}

func (a *Api) CreateTopic(ctx context.Context, in CreateTopicIn) error {
	u := fmt.Sprintf("https://console.aiven.io/v1beta/project/%v/service/%v/topic", in.Project, in.Service)
	out := struct {
		Errors []struct {
			Status  int
			Message string
		}
		Message string
	}{}

	if err := a.client.Post(ctx, u, in, &out); err != nil {
		return errors.Wrapf(err, "unable to create topic, %v, for project:service, %v:%v", in.TopicName, in.Project, in.Service)
	}

	for _, e := range out.Errors {
		if e.Status == http.StatusConflict {
			return nil
		}
		return fmt.Errorf(e.Message)
	}

	return nil
}

type DeleteTopicIn struct {
	Project   string
	Service   string
	TopicName string
}

func (a *Api) DeleteTopic(ctx context.Context, in DeleteTopicIn) error {
	u := fmt.Sprintf("https://console.aiven.io/v1beta/project/%v/service/%v/topic/%v", in.Project, in.Service, in.TopicName)
	out := struct {
		Errors []struct {
			Status  int
			Message string
		}
		Message string
	}{}

	if err := a.client.Delete(ctx, u, nil, &out); err != nil {
		return errors.Wrapf(err, "unable to create topic, %v, for project:service, %v:%v", in.TopicName, in.Project, in.Service)
	}

	for _, e := range out.Errors {
		if e.Status == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf(e.Message)
	}

	return nil
}

type TopicInfoIn struct {
	Project   string
	Service   string
	TopicName string
}

type Error struct {
	Message  string `json:"message"`
	MoreInfo string `json:"more_info"`
	Status   int    `json:"status"`
}

type ConsumerGroupInfo struct {
	GroupName string `json:"group_name"`
	Offset    int64  `json:"offset"`
}

type PartitionInfo struct {
	ConsumerGroups []ConsumerGroupInfo `json:"consumer_groups"`
	EarliestOffset int64               `json:"earliest_offset"`
	InSyncReplicas int                 `json:"isr"`
	LatestOffset   int64               `json:"latest_offset"`
	Partition      int32               `json:"partition"`
	Size           int64               `json:"size"`
}

type TopicInfo struct {
	CleanupPolicy     string          `json:"cleanup_policy"`
	MinInsyncReplicas int             `json:"min_insync_replicas"`
	Partitions        []PartitionInfo `json:"partitions"`
	Replication       int             `json:"replication"`
	RetentionBytes    int64           `json:"retention_bytes"`
	RetentionHours    int             `json:"retention_hours"`
	State             string          `json:"state"`
	TopicName         string          `json:"topic_name"`
}

type TopicInfoOut struct {
	Errors  []Error   `json:"errors"`
	Message string    `json:"message"`
	Topic   TopicInfo `json:"topic"`
}

// TopicInfo returns topic metadata
//
// See https://api.aiven.io/doc/#api-Service__Kafka-ServiceKafkaTopicGet
func (a *Api) TopicInfo(ctx context.Context, in TopicInfoIn) (TopicInfoOut, error) {
	out := TopicInfoOut{}
	u := fmt.Sprintf("https://console.aiven.io/v1beta/project/%v/service/%v/topic/%v", in.Project, in.Service, in.TopicName)
	if err := a.client.Get(ctx, u, &out); err != nil {
		return out, errors.Wrapf(err, "unable to retrieve topic info for topic, %v", in.TopicName)
	}

	return out, nil
}
