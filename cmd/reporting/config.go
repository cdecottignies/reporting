package main

import (
	"time"

	"gitlabdev.vadesecure.com/engineering/app/reporting/internal/jira"
	"gitlabdev.vadesecure.com/engineering/app/reporting/internal/storage"
	"gitlabdev.vadesecure.com/golib/u"
)

type Config struct {
	PublicAPI    APIConfig
	InternalAPI  APIConfig
	TechnicalAPI APIConfig
	Log          LogConfig
	Jira         jira.Config
	Storage      storage.DBConfig
}

type APIConfig struct {
	BindAddr     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type LogConfig struct {
	Level    string
	GinDebug bool
}

func ParseConfig() (*Config, error) {
	cfg := Config{
		PublicAPI: APIConfig{
			BindAddr:     "0.0.0.0:8080",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		InternalAPI: APIConfig{
			BindAddr:     "0.0.0.0:8081",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		TechnicalAPI: APIConfig{
			BindAddr:     "0.0.0.0:8082",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Log: LogConfig{
			Level:    "info",
			GinDebug: false,
		},
		Jira: jira.Config{
			URL:                "https://jira.cloudsystem.fr/",
			ConsumerKey:        "reporting_client",
			JiraPrivateKeyFile: `/conf/jira.pem`,
		},
		Storage: storage.DBConfig{
			Addr:           "mongodb:27017",
			MaxConnections: 10,
			DBName:         "reporting",
		},
	}

	err := u.DecodeFromEnv(&cfg, "REPORTING")
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
