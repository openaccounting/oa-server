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

## Docs

Api docs are available at https://openaccounting.io/api/

You may also build and run the docs locally with [apidoc](apidocjs.com) 

#### Requires [yarn](https://yarnpkg.com/) or [npm](https://www.npmjs.com/)

### Step 1:

Install `apidoc`

```bash
npm install apidoc -g
```

### Step 2:

Simply run `apidoc` within the source code root directory, this will automatically generate the documentation and write the output to`./doc` 

```bash
apidoc
```

### Step 3:

You may now navigate to the `./docs` directory and run any http server, or simply open `index.html` in your favorite browser!

## Help

[Join our Slack chatroom](https://join.slack.com/t/openaccounting/shared_invite/enQtNDc3NTAyNjYyOTYzLTc0ZjRjMzlhOTg5MmYwNGQxZGQyM2IzZTExZWE0NDFlODRlNGVhZmZiNDkyZDlhODYwZDcyNTQ5ZWJkMDU3N2M) and talk with us!
