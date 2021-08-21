test:
	go test -cover -race ./...

install:
	go install github.com/gopherguides/hype/cmd/hype