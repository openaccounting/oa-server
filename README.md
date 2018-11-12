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

## Help

[Join our Slack chatroom](https://openaccounting.slack.com/signup) and talk with us!