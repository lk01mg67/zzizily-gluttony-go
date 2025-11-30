# Gluttony

[![Go Report Card](https://goreportcard.com/badge/github.com/deuxksy/template-go-application)](https://goreportcard.com/report/github.com/deuxksy/template-go-application) Go Application 을 만들기 위한 기본 template

## Folder Layout

```bash
.
├── README.md
├── assets
├── build
│   └── ci
│       └── build.jenkinsfiles
├── cmd
│   └── template
│       └── main.go
├── configs
│   ├── dev.yml
│   └── local.yml
├── go.mod
├── go.sum
├── internal
│   ├── configuration
│   │   └── config_model.go
│   └── logger
│       └── logger.go
├── logs
│   └── 220707
│       ├── error.log
│       └── out.log
├── pkg
└── test
```

## Setup

1. Copy `.env.example` to `.env` (or set environment variables manually).
   ```bash
   cp .env.example .env
   ```
2. Edit `.env` with your credentials:
   ```bash
   USERID=your_id
   USERPW=your_password
   ```

## Configuration

Edit `configs/local.yml` (or create a new profile config) to define scenarios.
Example:
```yaml
Scenario:
  - name: Login
    url: https://assist9.i-on.net/login
    type: login
  - name: Healthcare
    url: ...
    type: booking
```

## mod & build

### windows

```bash
go mod tidy
GOOS=windows GOARCH=amd64 go build -o gluttony.exe cmd/gluttony/main.go
```

### linux

```bash
go mod tidy
GOOS=linux GOARCH=386 go build -o gluttony cmd/gluttony/main.go
```

### mac

```bash
go mod tidy
GOOS=darwin GOARCH=arm64 go build -o gluttony cmd/gluttony/main.go
```
