FROM golang:1.16-stretch

RUN apt update && apt install -y libasound2-dev python3

WORKDIR /earthly

deps:
    RUN go get golang.org/x/tools/cmd/goimports
    RUN go get golang.org/x/lint/golint
    RUN go get github.com/gordonklaus/ineffassign
    RUN go get github.com/jackc/tern
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

code:
    FROM +deps
    COPY --dir cmd ./

lint:
    FROM +code
    RUN output="$(ineffassign . | grep -v '/earthly/earthfile2llb/parser/.*\.go')" ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi
    RUN output="$(goimports -d . 2>&1)" | grep -v '.*\.pb\.go' ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi
    RUN golint -set_exit_status ./...
    RUN output="$(go vet ./... 2>&1)" ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi

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
