default: test install hype

test:
	go test -count 1 -race -vet=off -cover $$(go list ./... | grep -v /docs/)

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

hypecov:
	go test -coverprofile=coverage.out .
	go tool cover -html=coverage.out

install:
	go install -v ./cmd/hype

hype:
	hype export -format=markdown -f hype.md -o ./README.md

.PHONY: build docs
build:
	go build -o hype ./cmd/hype/

docs: build
	./hype export -f hype.md -format markdown -o README.md
