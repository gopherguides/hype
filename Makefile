default: test install

test:
	go test -timeout 10s -race -cover ./...

install:
	go install -v ./cmd/hype