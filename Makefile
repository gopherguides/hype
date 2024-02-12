default: test install hype

test:
	go test -timeout 10s -race -cover ./...

install:
	go install -v ./cmd/hype

hype:
	pushd .hype;hype export -format=markdown -f module.md > ../README.md;popd
