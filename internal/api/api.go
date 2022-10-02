package api

import (
	"context"
	"net/http"
	"time"

	"gitlabdev.vadesecure.com/engineering/app/reporting/internal/jira"
	"gitlabdev.vadesecure.com/engineering/app/reporting/internal/metrics"
	"gitlabdev.vadesecure.com/engineering/app/reporting/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"

	response "gitlabdev.vadesecure.com/engineering/app/reporting/client"
)

func Hello(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics.HelloCounter.WithLabelValues(scope).Inc()
		c.JSON(200, response.MsgResponse{Message: "Hello world " + scope})
	}
}

// IssueJira find issue by key in jira
func IssueJira(scope string, jiraC *jira.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		issue, codeR, err := jiraC.GetIssue(c.Param("key"))
		if err != nil {
			c.JSON(codeR, response.MsgResponse{Message: "Error :" + c.Param("key") + ": " + err.Error()})
		} else {
			c.JSON(codeR, response.IssueResponse{Project: issue.Project, Key: issue.Key, Type: issue.Type, Desc: issue.Desc, Assigned: issue.Assigned, Date: issue.Date})
		}
	}
}

// IssueMongo find issue by key in collection "issues"
func IssueMongo(scope string, mongoC *storage.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Param("key")) < 10 && len(c.Param("key")) > 4 {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			issue, err := mongoC.Issues.FindByKey(ctx, c.Param("key"))
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.JSON(400, response.MsgResponse{Message: "no issue with: " + c.Param("key") + ": " + err.Error()})
				} else {
					c.JSON(500, response.MsgResponse{Message: "Error: " + err.Error()})
				}
			} else {
				c.JSON(200, response.IssueResponse{Project: issue.Project, Key: issue.Key, Type: issue.Type, Desc: issue.Desc, Assigned: issue.Assigned, Date: issue.Date})
			}
		} else {
			c.JSON(400, response.MsgResponse{Message: "bad Param: " + c.Param("key")})

		}
	}
}

// Add add issue jira in the collection "issues"
func Add(scope string, mongoC *storage.Db, jiraC *jira.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := mongoC.Issues.FindByKey(ctx, c.Param("key"))
		if err != nil {
			if err == mongo.ErrNoDocuments {
				var issue *jira.Issue
				var codeR int
				issue, codeR, err = jiraC.GetIssue(c.Param("key"))
				if err != nil {
					c.JSON(codeR, response.MsgResponse{Message: "Error: " + c.Param("key") + ": " + err.Error()})
				} else {
					err = mongoC.Issues.Add(ctx, issue.Project, issue.Key, issue.Type, issue.Desc, issue.Assigned)
					if err != nil {
						c.JSON(500, response.MsgResponse{Message: "Error internal: " + err.Error()})
					}
					c.Status(codeR)
				}
			} else {
				c.JSON(500, response.MsgResponse{Message: "Error internal: " + c.Param("key") + "  " + err.Error()})
			}
		} else {
			c.JSON(400, response.MsgResponse{Message: "already added :" + c.Param("key")})
		}
	}
}

// Delete delete one issue in the collection "issues"
func Delete(scope string, mongoC *storage.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Param("key")) < 10 && len(c.Param("key")) > 4 {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, err := mongoC.Issues.FindByKey(ctx, c.Param("key"))
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.JSON(400, response.MsgResponse{Message: "no issue : " + c.Param("key") + ": " + err.Error()})
				} else {
					c.JSON(500, response.MsgResponse{Message: "error internal: " + c.Param("key") + ": " + err.Error()})
				}
			} else {
				err = mongoC.Issues.Delete(ctx, c.Param("key"))
				if err != nil {
					c.JSON(500, response.MsgResponse{Message: " " + err.Error()})
				}
				c.Status(http.StatusOK)
			}
		} else {
			c.JSON(400, response.MsgResponse{Message: "bad Param: " + c.Param("key")})
		}
	}
}

// All return all issue in the collection "issues"
func All(scope string, mongoC *storage.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		issue, err := mongoC.Issues.Find(ctx)
		if err != nil {
			c.JSON(500, response.MsgResponse{Message: "Error internal: " + err.Error()})
		} else {
			var t storage.Issue
			var IssuesRes []response.IssueResponse
			for i := 0; i < len(issue); i++ {
				t = issue[i]
				IssuesRes = append(IssuesRes, response.IssueResponse{Project: t.Project, Key: t.Key, Type: t.Type, Desc: t.Desc, Assigned: t.Assigned, Date: t.Date})
			}
			c.JSON(200, response.IssuesResponse{Issues: IssuesRes[:]})

		}
	}
}

// History return all issue in the collection "issues"
func History(scope string, mongoC *storage.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		issue, err := mongoC.Issues.History(ctx)
		if err != nil {
			c.JSON(500, response.MsgResponse{Message: "Error internal: " + err.Error()})
		} else {
			var t storage.Issue
			var IssuesRes []response.IssueResponse
			for i := 0; i < len(issue); i++ {
				t = issue[i]
				IssuesRes = append(IssuesRes, response.IssueResponse{Project: t.Project, Key: t.Key, Type: t.Type, Desc: t.Desc, Assigned: t.Assigned, Date: t.Date})
			}
			c.JSON(200, response.IssuesResponse{Issues: IssuesRes[:]})

		}
	}
}

// UpdateHistory copy the collection "issues" to the collection "history"
func UpdateHistory(scope string, mongoC *storage.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := mongoC.Issues.Updatehistory(ctx)
		if err != nil {
			c.JSON(500, response.MsgResponse{Message: "Error internal: " + err.Error()})
		} else {
			c.Status(http.StatusOK)
		}
	}
}

// Reset drop the collection
func Reset(scope string, mongoC *storage.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := mongoC.Issues.Reset(ctx)
		if err != nil {
			c.JSON(500, response.MsgResponse{Message: "Error internal: " + err.Error()})
		} else {
			c.Status(http.StatusOK)
		}
	}
}
