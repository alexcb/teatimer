FROM golang:1.13-alpine3.11

RUN apk add --update --no-cache \
    alsa-lib-dev \
    bash \
    bash-completion \
    binutils \
    ca-certificates \
    coreutils \
    curl \
    findutils \
    g++ \
    git \
    gnupg \
    grep \
    less \
    make \
    openssl \
    protoc \
    shellcheck \
    util-linux

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
    ARG EARTHLY_GIT_HASH
    RUN go build -ldflags  "-X main.GitSha=$EARTHLY_GIT_HASH" -o teatimer cmd/teatimer/main.go
    SAVE ARTIFACT teatimer AS LOCAL .

all:
    BUILD +lint
    BUILD +teatimer
