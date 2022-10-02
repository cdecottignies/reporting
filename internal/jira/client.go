package jira

import (
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/andygrunwald/go-jira"
)

type Config struct {
	URL                string
	ConsumerKey        string
	JiraPrivateKeyFile string
}
type Client struct {
	jiraClient *jira.Client
}
type Issue struct {
	Project string `bson:"project"`
	// Key unique key format : VTM-84
	Key      string `bson:"key"`
	Type     string `bson:"type"`
	Desc     string `bson:"desc"`
	Assigned string `bson:"assigned"`
	//Date format : yyyy-n°week = 2021-18
	Date string `bson:"date"`
}

// NewClient new client and need privatekeyjira
func NewClient(cfg *Config) (*Client, error) {
	var c Client
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(cfg.JiraPrivateKeyFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := make([]byte, 4096)
	var n int
	n, err = f.Read(buf)
	if err != nil {
		return nil, err
	}
	privateKey := string(buf[0:n])
	c.jiraClient, err = getJIRAClient(u, cfg.ConsumerKey, privateKey)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// formyearWeek Date format : yyyy-n°week = 2021-18
func formyearWeek() string {
	now := time.Now().UTC()
	year, week := now.ISOWeek()
	var date string = strconv.Itoa(year) + "-" + strconv.Itoa(week)
	return date

}

// GetIssue whith key jira for return one struct Issue
func (c *Client) GetIssue(key string) (*Issue, int, error) {
	issue, response, err := c.jiraClient.Issue.Get(key, nil)
	res := response.StatusCode
	if err != nil {
		return nil, res, err
	}
	var result Issue
	result.Project = issue.Fields.Project.Name
	result.Key = issue.Key
	result.Desc = issue.Fields.Summary
	result.Type = issue.Fields.Type.Name
	result.Date = formyearWeek()
	if issue.Fields.Assignee != nil {
		result.Assigned = issue.Fields.Assignee.Name
	}
	return &result, res, err
}
