build-api:
	go build cmd/api-server/api-server.go

build-crawler:
	go build cmd/crawler-server/crawler-server.go

.PHONY: build-api build-crawler