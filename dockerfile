ARG APP_NAME=teminel
FROM golang:1.21-alpine as builder
ARG APP_NAME
ENV CGO_ENABLED=0
ARG GOOS=linux
ENV GOOS=${GOOS}
ARG GOARCH=amd64
ENV GOARCH=${GOARCH}

RUN apk add upx

WORKDIR /usr/src/${APP_NAME}

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -ldflags="-s -w" -o .build/${APP_NAME} .
RUN upx --brute .build/${APP_NAME}

FROM scratch
ARG APP_NAME
COPY --from=builder /usr/src/${APP_NAME}/.build/${APP_NAME} /usr/bin/${APP_NAME}
ENTRYPOINT ["${APP_NAME}"]
