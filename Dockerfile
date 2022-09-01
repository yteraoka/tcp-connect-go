FROM golang:1.14.15-alpine3.11 as builder

WORKDIR /work

COPY go.* *.go ./

RUN go mod download && CGO_ENABLED=0 GOOS=linux go build

FROM golang:1.14.15-alpine3.11
WORKDIR /
COPY --from=builder /work/tcp-connect-go /usr/bin/tcp-connect-go
