default: test install hype

test:
	go test -timeout 10s -race -cover ./...

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

hypecov:
	go test -coverprofile=coverage.out .
	go tool cover -html=coverage.out

install:
	go install -v ./cmd/hype

hype:
	pushd .hype;hype export -format=markdown -f hype.md > ../README.md;popd
