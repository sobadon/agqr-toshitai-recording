.PHONY: build
build:
	go build -trimpath ./cmd/agqr-toshitai-recording

.PHONY: test
test:
	go test -v ./...
