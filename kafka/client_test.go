package kafka_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/savaki/aiven"
	"github.com/savaki/aiven/kafka"
	"github.com/tj/assert"
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

	api := kafka.New(client)
	out, err := api.TopicInfo(context.Background(), kafka.TopicInfoIn{
		Project:   project,
		Service:   service,
		TopicName: topic,
	})
	assert.Nil(t, err)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(out)
}
