package aiven

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Kafka provides an api into aiven kafka
type Kafka struct {
	client *Client
}

// New accepts a valid aiven client and returns access to the Kafka api
func newKafka(client *Client) *Kafka {
	return &Kafka{
		client: client,
	}
}

// KafkaTopic represents the Kafka topics
type KafkaTopic struct {
	CleanupPolicy  string `json:"cleanup_policy"`
	Partitions     int    `json:"partitions"`
	Replication    int    `json:"replication"`
	RetentionHours int    `json:"retention_hours"`
	State          string `json:"state"`
	TopicName      string `json:"topic_name"`
}

type KafkaListTopicsIn struct {
	Project string
	Service string
}

// ListTopics returns the list of all topics
func (k *Kafka) ListTopics(ctx context.Context, in KafkaListTopicsIn) ([]KafkaTopic, error) {
	u := fmt.Sprintf("https://console.aiven.io/v1beta/project/%v/service/%v", in.Project, in.Service)
	out := struct {
		Service struct {
			Topics []KafkaTopic
		}
	}{}
	if err := k.client.Get(ctx, u, &out); err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve topics for project:service, %v:%v", in.Project, in.Service)
	}

	return out.Service.Topics, nil
}

type KafkaCreateTopicIn struct {
	Project        string `json:"-"`
	Service        string `json:"-"`
	CleanupPolicy  string `json:"cleanup_policy"`
	Partitions     int    `json:"partitions"`
	Replication    int    `json:"replication"`
	RetentionHours int    `json:"retention_hours"`
	TopicName      string `json:"topic_name"`
}

func (k *Kafka) CreateTopic(ctx context.Context, in KafkaCreateTopicIn) error {
	u := fmt.Sprintf("https://console.aiven.io/v1beta/project/%v/service/%v/topic", in.Project, in.Service)
	out := struct {
		Errors []struct {
			Status  int
			Message string
		}
		Message string
	}{}

	if err := k.client.Post(ctx, u, in, &out); err != nil {
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

type KafkaDeleteTopicIn struct {
	Project   string
	Service   string
	TopicName string
}

func (k *Kafka) DeleteTopic(ctx context.Context, in KafkaDeleteTopicIn) error {
	u := fmt.Sprintf("https://console.aiven.io/v1beta/project/%v/service/%v/topic/%v", in.Project, in.Service, in.TopicName)
	out := struct {
		Errors []struct {
			Status  int
			Message string
		}
		Message string
	}{}

	if err := k.client.Delete(ctx, u, nil, &out); err != nil {
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

type KafkaTopicInfoIn struct {
	Project   string
	Service   string
	TopicName string
}

type Error struct {
	Message  string `json:"message"`
	MoreInfo string `json:"more_info"`
	Status   int    `json:"status"`
}

type KafkaConsumerGroupInfo struct {
	GroupName string `json:"group_name"`
	Offset    int64  `json:"offset"`
}

type KafkaPartitionInfo struct {
	ConsumerGroups []KafkaConsumerGroupInfo `json:"consumer_groups"`
	EarliestOffset int64                    `json:"earliest_offset"`
	InSyncReplicas int                      `json:"isr"`
	LatestOffset   int64                    `json:"latest_offset"`
	Partition      int32                    `json:"partition"`
	Size           int64                    `json:"size"`
}

type KafkaTopicInfo struct {
	CleanupPolicy     string               `json:"cleanup_policy"`
	MinInsyncReplicas int                  `json:"min_insync_replicas"`
	Partitions        []KafkaPartitionInfo `json:"partitions"`
	Replication       int                  `json:"replication"`
	RetentionBytes    int64                `json:"retention_bytes"`
	RetentionHours    int                  `json:"retention_hours"`
	State             string               `json:"state"`
	TopicName         string               `json:"topic_name"`
}

type KafkaTopicInfoOut struct {
	Errors  []Error        `json:"errors"`
	Message string         `json:"message"`
	Topic   KafkaTopicInfo `json:"topic"`
}

// KafkaTopicInfo returns topic metadata
//
// See https://api.aiven.io/doc/#api-Service__Kafka-ServiceKafkaTopicGet
func (k *Kafka) TopicInfo(ctx context.Context, in KafkaTopicInfoIn) (KafkaTopicInfoOut, error) {
	out := KafkaTopicInfoOut{}
	u := fmt.Sprintf("https://console.aiven.io/v1beta/project/%v/service/%v/topic/%v", in.Project, in.Service, in.TopicName)
	if err := k.client.Get(ctx, u, &out); err != nil {
		return out, errors.Wrapf(err, "unable to retrieve topic info for topic, %v", in.TopicName)
	}

	return out, nil
}
