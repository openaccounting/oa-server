# Open Accounting Server

## Prerequisites

1. Go 1.8+
2. MySQL 5.7+

## Database setup

Use schema.sql and indexes.sql to create a MySQL database to store Open Accounting data.

## Configuration

Copy config.json.sample to config.json and edit to match your information.

## Run

`go run core/server.go`

## Build

`go build core/server.go`

## Docker

If you are interested in running Open Accounting via Docker, @alokmenghrajani has created a [repo](https://github.com/alokmenghrajani/openaccounting-docker) for this.

## Help

[Join our Slack chatroom](https://join.slack.com/t/openaccounting/shared_invite/zt-23zy988e8-93HP1GfLDB7osoQ6umpfiA) and talk with us!