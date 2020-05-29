.PHONY: all deps test download-run download-build upload-run upload-build

all: test

deps:
	@docker run --rm -it -v "${GOPATH}":/go \
			-w /go/src/github.com/pottava/gcs-utils \
			golang:1.14.3-alpine3.11 go mod vendor

test: deps
	@docker run --rm -it -v "${GOPATH}":/go \
			-w /go/src/github.com/pottava/gcs-utils \
			golangci/golangci-lint:v1.27.0-alpine \
			golangci-lint run --config .golangci.yml
	@docker run --rm -it -v "${GOPATH}":/go \
			-w /go/src/github.com/pottava/gcs-utils \
			golangci/golangci-lint:v1.27.0-alpine \
			sh -c "go list ./... | grep -v /vendor/ | xargs go test -p 1 -count=1"

download-run:
	@docker run --rm -it -v "${GOPATH}":/go \
			-v "${HOME}/.config/gcloud":/root/.config/gcloud \
			-w /go/src/github.com/pottava/gcs-utils \
			golang:1.14.3-alpine3.11 \
			go run cmd/download/main.go

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

upload-run:
	@docker run --rm -it -v "${GOPATH}":/go \
			-v "${HOME}/.config/gcloud":/root/.config/gcloud \
			-w /go/src/github.com/pottava/gcs-utils \
			golang:1.14.3-alpine3.11 \
			go run cmd/upload/main.go

upload-build:
	@docker run --rm -it -v "${GOPATH}":/go \
			-w /go/src/github.com/pottava/gcs-utils/cmd/upload \
			-e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 \
			golang:1.14.3-alpine3.11 \
			go build -ldflags "-s -w" -o "dist/gcs-upload-linux"
	@docker run --rm -it -v "${GOPATH}":/go \
			-w /go/src/github.com/pottava/gcs-utils/cmd/upload \
			-e GOOS=darwin -e GOARCH=amd64 -e CGO_ENABLED=0 \
			golang:1.14.3-alpine3.11 \
			go build -ldflags "-s -w" -o "dist/gcs-upload-mac"
	@docker run --rm -it -v "${GOPATH}":/go \
			-w /go/src/github.com/pottava/gcs-utils/cmd/upload \
			-e GOOS=windows -e GOARCH=amd64 -e CGO_ENABLED=0 \
			golang:1.14.3-alpine3.11 \
			go build -ldflags "-s -w" -o "dist/gcs-upload.exe"
