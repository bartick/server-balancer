# builder stage
FROM golang:1.21.0-bookworm AS builder

WORKDIR /server

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -o app

# final stage
FROM ubuntu:22.04

WORKDIR /go

COPY --from=builder server/app .

CMD ["./app"]
