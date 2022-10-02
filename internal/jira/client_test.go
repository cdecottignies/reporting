package jira

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDb(t *testing.T) {
	client, err := NewClient(
		&Config{
			URL:                "https://jira.cloudsystem.fr/",
			ConsumerKey:        "reporting_client",
			JiraPrivateKeyFile: `/conf/jira.pem`,
		})
	if err != nil {
		t.Skipf("disabled, no connection to api jira")
		return
	}
	require.Nil(t, err)
	require.NotNil(t, client)

	t.Run("formyearWeek", func(t *testing.T) {
		res := formyearWeek()
		require.NotNil(t, res)
		now := time.Now().UTC()
		year, week := now.ISOWeek()
		assert.Equal(t, res, strconv.Itoa(year)+"-"+strconv.Itoa(week))

	})
	t.Run("get issue", func(t *testing.T) {
		res, response, err := client.GetIssue("ENG-6738")
		require.Nil(t, err)
		require.NotNil(t, res)
		assert.Equal(t, response, 200)
		assert.Equal(t, res.Assigned, "cdecottignies")
		assert.Equal(t, res.Key, "ENG-6738")
		assert.Equal(t, res.Type, "Task")
		assert.Equal(t, res.Project, "Engineering")

	})
	t.Run("fake get issue", func(t *testing.T) {
		res, response, err := client.GetIssue("ENG-9999")
		assert.NotEqual(t, response, 200)

		require.NotNil(t, err)
		require.Nil(t, res)

	})
}
