VERSION 0.7
FROM golang:1.20-bookworm

RUN apt update && apt install -y libasound2-dev python3 curl

WORKDIR /earthly

deps:
    RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.0
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

code:
    FROM +deps
    COPY --dir cmd ./

lint:
    FROM +code
    COPY ./.golangci.yaml ./
    RUN golangci-lint run

teatimer:
    FROM +code
    COPY sounds/ready.mp3 cmd/teatimer/. # sound taken from https://github.com/faiface/beep
    ARG EARTHLY_GIT_HASH
    ENV GOOS=linux
    ENV GOARCH=amd64
    RUN go build -ldflags  "-X main.GitSha=$EARTHLY_GIT_HASH" -o teatimer cmd/teatimer/main.go
    SAVE ARTIFACT teatimer AS LOCAL teatimer

all:
    BUILD +lint
    BUILD +teatimer
