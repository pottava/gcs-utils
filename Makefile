.PHONY: download-run download-deps download-test download-build

download-run:
	@docker run --rm -it -v "${GOPATH}":/go \
			-v "${HOME}/.config/gcloud":/root/.config/gcloud \
			-w /go/src/github.com/pottava/gcs-utils \
			golang:1.14.3-alpine3.11 \
			go run cmd/download/main.go

download-deps:
	@docker run --rm -it -v "${GOPATH}":/go \
			-w /go/src/github.com/pottava/gcs-utils \
			golang:1.14.3-alpine3.11 go mod vendor

download-test:
	@docker run --rm -it -v "${GOPATH}":/go \
			-w /go/src/github.com/pottava/gcs-utils \
			golangci/golangci-lint:v1.27.0-alpine \
			golangci-lint run --config .golangci.yml
	@docker run --rm -it -v "${GOPATH}":/go \
			-w /go/src/github.com/pottava/gcs-utils \
			golangci/golangci-lint:v1.27.0-alpine \
			sh -c "go list ./... | grep -v /vendor/ | xargs go test -p 1 -count=1"

download-build:
	@docker run --rm -it -v "${GOPATH}":/go \
			-w /go/src/github.com/pottava/gcs-utils/cmd/download \
			-e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 \
			golang:1.14.3-alpine3.11 \
			go build -ldflags "-s -w" -o "dist/gcs-download-linux"
	@docker run --rm -it -v "${GOPATH}":/go \
			-w /go/src/github.com/pottava/gcs-utils/cmd/download \
			-e GOOS=darwin -e GOARCH=amd64 -e CGO_ENABLED=0 \
			golang:1.14.3-alpine3.11 \
			go build -ldflags "-s -w" -o "dist/gcs-download-mac"
	@docker run --rm -it -v "${GOPATH}":/go \
			-w /go/src/github.com/pottava/gcs-utils/cmd/download \
			-e GOOS=windows -e GOARCH=amd64 -e CGO_ENABLED=0 \
			golang:1.14.3-alpine3.11 \
			go build -ldflags "-s -w" -o "dist/gcs-download.exe"
