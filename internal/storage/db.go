package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBConfig is a MongoDB configuration
type DBConfig struct {
	// Addr is the address of the server
	Addr string
	// MaxConnections is the max of connections
	MaxConnections int
	// DBName is the database name
	DBName string
}

type Db struct {
	Issues Issues
}

func New(cfg *DBConfig) (*Db, error) {
	var res Db

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", cfg.Addr))
	if cfg.MaxConnections > 0 {
		clientOptions.SetMaxPoolSize(uint64(cfg.MaxConnections))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	res.Issues.c = client.Database(cfg.DBName).Collection("issues")

	return &res, nil

}
