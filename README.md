# Reporting HTTP API service


This service uses the [Gin Web Framework](https://github.com/gin-gonic/gin).

## `cmd` directory

This directory contains all the program. For HTTP API services, it contains at least
the server.

## `internal` directory

### `api` package

GET:<br>
-get IssueJira<br>
-get IssueMongo<br>
-Add<br>
-Delete<br>
-All<br>
-History<br>
-UpdateHistory<br>
-Reset<br>

### `jira` package

This package use the [go-jira librairy](https://pkg.go.dev/github.com/andygrunwald/go-jira@v1.13.0?utm_source=gopls).<br>
Oauth1 of the API jira, use the directory 'conf' for the token jira and private key<br>
only one request 'Get jira issue'.<br>

### `storage` package

This package uses the [Go mongo-driver](https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.5.2).<br>
mongo driver connect to the image mongoDB.<br>
there classic MongoDB requests used:<br>

-Add<br>
-AddMany<br>
-Delete<br>
-Find<br>
-FindOne<br>

### `metrics` package

This package contains all the exposed metrics.

## `tests` directory

This directory contains all the functional tests written in python using the
[pyvade framework](https://gitlabdev.vadesecure.com/engineering/pyvade)

## `client` package

The `client` package contains the go client of the API.

## `ui` directory

routers:<br>
    path: '/',name: 'Home'<br>
    path: '/user',name: 'user'<br>
    path: '/admin',name: 'admin'<br>


## `doc` directory

This directory contains all the documentation about this service (configuration, OpenAPI specification...).
