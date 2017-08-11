package aiven_test

import (
	"os"
	"testing"

	"github.com/savaki/aiven"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	email := os.Getenv("AIVEN_EMAIL")
	password := os.Getenv("AIVEN_PASSWORD")
	project := os.Getenv("AIVEN_PROJECT")
	service := os.Getenv("AIVEN_SERVICE")

	if email == "" {
		t.Skip("AIVEN_EMAIL not set")
		return
	}
	if password == "" {
		t.Skip("AIVEN_PASSWORD not set")
		return
	}
	if project == "" {
		t.Skip("AIVEN_PROJECT not set")
		return
	}
	if service == "" {
		t.Skip("AIVEN_SERVICE not set")
		return
	}

	client, err := aiven.New(email, password)
	assert.Nil(t, err)
	assert.NotNil(t, client)
}
