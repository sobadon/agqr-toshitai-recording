.PHONY: build
build:
	CGO_ENABLED=0 go build -trimpath ./cmd/agqr-toshitai-recording

.PHONY: test
test:
	go test -v ./...
