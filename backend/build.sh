#!/usr/bin/env bash

GOOS=linux GOARCH=amd64 go build -o recover_videos cmd/recover_videos/main.go
zip recover_videos.zip client_secret.json recover_videos .env
rm recover_videos

GOOS=linux GOARCH=amd64 go build -o hard_migrate cmd/hard_migrate/main.go
zip hard_migrate.zip hard_migrate .env
rm hard_migrate

GOOS=linux GOARCH=amd64 go build -o soft_migrate cmd/soft_migrate/main.go
zip soft_migrate.zip soft_migrate .env
rm soft_migrate
