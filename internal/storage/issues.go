package storage

import (
	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Issues struct {
	c *mongo.Collection
}

func (i *Issues) Reset(ctx context.Context) error {
	return i.c.Drop(ctx)
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

// FormyearWeek create date format yyyy-n°week
func FormyearWeek() string {
	now := time.Now().UTC()
	year, week := now.ISOWeek()
	var date string = strconv.Itoa(year) + "-" + strconv.Itoa(week)
	return date

}

// AddIssue add one issue with Issue
func (i *Issues) AddIssue(ctx context.Context, arg Issue) error {
	if arg.Date == "" {
		arg.Date = FormyearWeek()
	}
	_, err := i.c.InsertOne(ctx, arg)

	return err
}

// Add add one issue with year and n°week
func (i *Issues) Add(ctx context.Context, Project string, Key string, Type string, Desc string, Assigned string) error {
	res := FormyearWeek()
	doc := Issue{Project: Project, Key: Key, Type: Type, Desc: Desc, Assigned: Assigned, Date: res}
	_, err := i.c.InsertOne(ctx, doc)
	return err
}

// AddMany add many issue
func (i *Issues) AddMany(ctx context.Context, arg []Issue) error {
	res := []interface{}{}
	for i := 0; i < len(arg); i++ {
		res = append(res, Issue{Project: arg[i].Project, Key: arg[i].Key, Type: arg[i].Type, Desc: arg[i].Desc, Assigned: arg[i].Assigned, Date: arg[i].Date})
	}
	_, err := i.c.InsertMany(ctx, res)
	if err != nil {
		return err
	}
	return nil
}

// Delete delete one issue by key
func (i *Issues) Delete(ctx context.Context, Key string) error {
	doc := bson.M{"key": Key}
	_, err := i.c.DeleteOne(ctx, doc)
	return err
}

// FindByKey find one issue by key (only 1 assigned because decode)
func (i *Issues) FindByKey(ctx context.Context, Key string) (*Issue, error) {
	doc := bson.M{"key": Key}
	var result Issue
	err := i.c.FindOne(ctx, doc).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindByProject find all issues of a project
func (i *Issues) FindByProject(ctx context.Context, Project string) ([]Issue, error) {
	doc := bson.M{"project": Project}
	var result []Issue
	cur, err := i.c.Find(ctx, doc)
	for i := 0; cur.Next(ctx); i++ {
		var e Issue
		err := cur.Decode(&e)
		if err != nil {
			return nil, err
		}
		result = append(result, e)

	}
	return result, err
}

// Find find the collection of issue
func (i *Issues) Find(ctx context.Context) ([]Issue, error) {
	doc := bson.M{}
	var result []Issue
	cur, err := i.c.Find(ctx, doc)
	for i := 0; cur.Next(ctx); i++ {
		var e Issue
		err := cur.Decode(&e)
		if err != nil {
			return nil, err
		}
		result = append(result, e)

	}
	return result, err
}

// History return all the collection history
func (i *Issues) History(ctx context.Context) ([]Issue, error) {
	i.c = i.c.Database().Collection("history")
	issues, err := i.Find(ctx)
	if err != nil {
		i.c = i.c.Database().Collection("issues")
		return nil, err
	}
	i.c = i.c.Database().Collection("issues")
	return issues, err
}

// Updatehistory copy the collection issues to the collection history
func (i *Issues) Updatehistory(ctx context.Context) error {
	var results []Issue
	results, err := i.Find(ctx)
	if err != nil {
		return err
	}
	i.c = i.c.Database().Collection("history")
	err = i.AddMany(ctx, results)
	i.c = i.c.Database().Collection("issues")
	return err
}
