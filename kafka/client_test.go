package kafka_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/savaki/aiven"
	"github.com/savaki/aiven/kafka"
	"github.com/stretchr/testify/assert"
)

func TestTopics(t *testing.T) {
	project := os.Getenv("AIVEN_PROJECT")
	service := os.Getenv("AIVEN_SERVICE")

	if project == "" {
		t.Skip("AIVEN_PROJECT not set")
		return
	}
	if service == "" {
		t.Skip("AIVEN_SERVICE not set")
		return
	}

	client, err := aiven.EnvAuth()
	assert.Nil(t, err)
	assert.NotNil(t, client)

	api := kafka.New(client)
	fmt.Println(api.Topics(context.Background(), project, service))
}
