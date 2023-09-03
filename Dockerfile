# builder stage
FROM --platform=$BUILDPLATFORM golang:bullseye AS builder

WORKDIR /server

COPY . .

RUN go mod download

ENV GOOS=linux

ENV GOARCH=amd64

RUN go build -o app

# final stage
FROM ubuntu:22.04

WORKDIR /go

COPY --from=builder server .

ENV ENVIRONMENT=production

CMD ["./app"]
