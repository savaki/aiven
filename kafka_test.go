package aiven_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/savaki/aiven"
	"github.com/stretchr/testify/assert"
)

func TestTopics(t *testing.T) {
	project := os.Getenv("AIVEN_PROJECT")
	service := os.Getenv("AIVEN_SERVICE")
	topic := os.Getenv("AIVEN_TOPIC")

	if project == "" {
		t.Skip("AIVEN_PROJECT not set")
	}
	if service == "" {
		t.Skip("AIVEN_SERVICE not set")
	}
	if topic == "" {
		t.Skip("AIVEN_TOPIC not set")
	}

	client, err := aiven.EnvAuth()
	assert.Nil(t, err)
	assert.NotNil(t, client)

	api := client.Kafka()
	out, err := api.TopicInfo(context.Background(), aiven.KafkaTopicInfoIn{
		Project:   project,
		Service:   service,
		TopicName: topic,
	})
	assert.Nil(t, err)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(out)
}
