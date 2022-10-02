package storage

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDb connect to the MongoDB server
func TestDb(t *testing.T) {
	db, err := New(
		&DBConfig{
			Addr:           "172.19.0.2:27017",
			MaxConnections: 10,
			DBName:         "test",
		})
	if err != nil {
		t.Skipf("disabled, no connection do real mongodb")
		return
	}
	require.NotNil(t, db)

	ctx := context.Background()
	db.Issues.Reset(ctx)
	t.Run("FindByKey", func(t *testing.T) {
		doc := Issue{Project: "heimdall", Key: "HEIM-88", Type: "bug", Desc: "error 500", Assigned: "clem", Date: "2021-19"}
		_, err = db.Issues.c.InsertOne(ctx, doc)
		require.Nil(t, err)
		res, err := db.Issues.FindByKey(ctx, "HEIM-88")
		require.Nil(t, err)
		assert.Equal(t, res.Project, "heimdall")
		assert.Equal(t, res.Key, "HEIM-88")
		assert.Equal(t, res.Desc, "error 500")
		err = db.Issues.Reset(ctx)
		require.Nil(t, err)

	})
	t.Run("AddIssue", func(t *testing.T) {
		err = db.Issues.AddIssue(ctx, Issue{Project: "heimdall", Key: "HEIM-87", Type: "bug", Desc: "error 404", Assigned: "clem", Date: "2021-18"})
		require.Nil(t, err)
		var res *Issue
		res, err := db.Issues.FindByKey(ctx, "HEIM-87")
		require.Nil(t, err)
		assert.Equal(t, res.Assigned, "clem")
		assert.Equal(t, res.Key, "HEIM-87")
		require.Nil(t, err)
		assert.Equal(t, res.Assigned, "clem")
		assert.Equal(t, res.Key, "HEIM-87")
		assert.Equal(t, res.Type, "bug")
		assert.Equal(t, res.Project, "heimdall")
		err = db.Issues.Reset(ctx)
		require.Nil(t, err)

	})
	t.Run("Add", func(t *testing.T) {
		err = db.Issues.Add(ctx, "heimdall", "HEIM-87", "bug", "error 404", "clem")
		require.Nil(t, err)
		var res *Issue
		res, err := db.Issues.FindByKey(ctx, "HEIM-87")
		require.Nil(t, err)

		assert.Equal(t, res.Assigned, "clem")
		assert.Equal(t, res.Key, "HEIM-87")
		assert.Equal(t, res.Type, "bug")
		assert.Equal(t, res.Project, "heimdall")
		now := time.Now().UTC()
		year, week := now.ISOWeek()
		var date string = strconv.Itoa(year) + "-" + strconv.Itoa(week)
		assert.Equal(t, res.Date, date)
		err = db.Issues.Reset(ctx)
		require.Nil(t, err)

	})
	t.Run("Delete", func(t *testing.T) {
		err = db.Issues.Add(ctx, "heimdall", "HEIM-87", "bug", "error 404", "clem")
		require.Nil(t, err)
		err = db.Issues.Delete(ctx, "HEIM-87")
		require.Nil(t, err)
		_, err := db.Issues.FindByKey(ctx, "HEIM-87")
		require.NotNil(t, err)
		assert.NotNil(t, err)
		err = db.Issues.Reset(ctx)
		require.Nil(t, err)

	})
	t.Run("FindByProject", func(t *testing.T) {
		err = db.Issues.Add(ctx, "heimdall", "HEIM-88", "bug", "error 500", "clem")
		require.Nil(t, err)
		err = db.Issues.Add(ctx, "heimdall", "HEIM-89", "bug", "error segmentation", "clem")
		require.Nil(t, err)
		var res []Issue
		res, err := db.Issues.FindByProject(ctx, "heimdall")
		require.Nil(t, err)

		assert.Equal(t, res[0].Project, "heimdall")
		assert.Equal(t, res[1].Project, "heimdall")
		assert.Equal(t, res[0].Key, "HEIM-88")
		assert.Equal(t, res[1].Key, "HEIM-89")
		assert.Equal(t, res[0].Desc, "error 500")
		assert.Equal(t, res[1].Desc, "error segmentation")
		err = db.Issues.Reset(ctx)
		require.Nil(t, err)

	})
	t.Run("Find", func(t *testing.T) {
		err = db.Issues.Add(ctx, "heimdall", "HEIM-88", "bug", "error 500", "clem")
		require.Nil(t, err)
		err = db.Issues.Add(ctx, "heimdall", "HEIM-89", "bug", "error segmentation", "clem")
		require.Nil(t, err)
		var res []Issue
		res, err := db.Issues.Find(ctx)
		require.Nil(t, err)

		assert.Equal(t, res[0].Project, "heimdall")
		assert.Equal(t, res[1].Project, "heimdall")
		assert.Equal(t, res[0].Key, "HEIM-88")
		assert.Equal(t, res[1].Key, "HEIM-89")
		assert.Equal(t, res[0].Desc, "error 500")
		assert.Equal(t, res[1].Desc, "error segmentation")
		err = db.Issues.Reset(ctx)
		require.Nil(t, err)

	})
	t.Run("AddMany", func(t *testing.T) {
		var tab []Issue
		doc := Issue{Project: "reporting", Key: "REP-10", Type: "task", Desc: "pipeline", Assigned: "clem", Date: "2021-20"}
		tab = append(tab, doc)
		doc = Issue{Project: "reporting", Key: "REP-11", Type: "bug", Desc: "docker", Assigned: "clem", Date: "2021-20"}
		tab = append(tab, doc)
		err = db.Issues.AddMany(ctx, tab)
		require.Nil(t, err)
		var res *Issue
		res, err := db.Issues.FindByKey(ctx, "REP-11")
		require.Nil(t, err)
		assert.Equal(t, res.Assigned, "clem")
		assert.Equal(t, res.Key, "REP-11")
		assert.Equal(t, res.Type, "bug")
		assert.Equal(t, res.Project, "reporting")
		err = db.Issues.Reset(ctx)
		require.Nil(t, err)
	})

	t.Run("Updatehistory", func(t *testing.T) {
		err = db.Issues.Add(ctx, "heimdall", "HEIM-88", "bug", "error 500", "clem")
		require.Nil(t, err)
		err = db.Issues.Add(ctx, "heimdall", "HEIM-89", "bug", "error", "clem")
		require.Nil(t, err)
		err = db.Issues.Updatehistory(ctx)
		require.Nil(t, err)
		err = db.Issues.Reset(ctx)
		require.Nil(t, err)
		db.Issues.c = db.Issues.c.Database().Collection("history")
		res, err := db.Issues.FindByKey(ctx, "HEIM-89")
		require.Nil(t, err)
		assert.Equal(t, res.Assigned, "clem")
		assert.Equal(t, res.Key, "HEIM-89")
		assert.Equal(t, res.Type, "bug")
		assert.Equal(t, res.Project, "heimdall")
		err = db.Issues.Reset(ctx)
		require.Nil(t, err)
		db.Issues.c = db.Issues.c.Database().Collection("issues")

	})

}
