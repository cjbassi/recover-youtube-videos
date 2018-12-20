.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/hard-migrate cmd/hard_migrate/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/soft-migrate cmd/soft_migrate/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/fetch-removed-videos cmd/fetch_removed_videos/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
