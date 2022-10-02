package client

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getRealClient(t *testing.T) *Client {
	exampleAddr := os.Getenv("TESTS_EXAMPLE_BASE_ADDR")
	if exampleAddr == "" {
		t.Skip("Missing environment variable TESTS_EXAMPLE_BASE_ADDR")
	}

	c, err := New(Config{
		BaseURL: exampleAddr,
	})
	assert.NoError(t, err)

	return c
}

// TestFunctionalPing must be run with a real example service with no data.
func TestFunctionalPing(t *testing.T) {
	client := getRealClient(t)

	err := client.Ping(context.Background())
	require.NoError(t, err)
}

// TestFunctionalEvents must be run with a real example service with no data.
func TestFunctionalHello(t *testing.T) {
	client := getRealClient(t)

	resp, err := client.Hello(context.Background())
	require.NoError(t, err)
	require.Equal(t, &MsgResponse{Message: "Hello world public"}, resp)
}
