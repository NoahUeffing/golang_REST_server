# syntax=docker/dockerfile:1
FROM golang:1.16-alpine
RUN apk add build-base

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY *.ico ./
COPY *.sqlite ./
RUN go build -o /golang-rest-server
EXPOSE 4041
CMD [ "/golang-rest-server" ]