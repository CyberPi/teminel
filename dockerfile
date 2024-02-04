ARG APP_NAME=teminel
FROM golang:1.21-alpine as builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -ldflags="-s -w" -o teminel ./...

CMD ["app"]
