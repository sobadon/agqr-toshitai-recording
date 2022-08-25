.PHONY: build
build:
	go build -trimpath ./cmd/agqr-toshitai-recording

.PHONY: test
test:
	go test -v ./...

.PHONY: gen-mock
gen-mock:
	rm -r ./testdata/mock
	go generate ./...
