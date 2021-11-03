test:
	go test -cover -race ./...

install: generate
	go install github.com/gopherguides/hype/cmd/hype

generate:
	go generate ./...