default: test install hype

test:
	go test -timeout 10s -count 1 -race -cover $$(go list ./... | grep -v /docs/)

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

hypecov:
	go test -coverprofile=coverage.out .
	go tool cover -html=coverage.out

install:
	go install -v ./cmd/hype

hype:
	hype export -format=markdown -f hype.md > ./README.md
