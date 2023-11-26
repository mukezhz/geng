FROM golang:alpine

# Required because go requires gcc to build
RUN apk add build-base git inotify-tools
RUN echo $GOPATH
RUN go install github.com/rubenv/sql-migrate/...@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /clean_web
COPY . .
RUN go mod download

CMD sh docker/run.sh
