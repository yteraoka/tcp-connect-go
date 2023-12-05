FROM golang:1.21.5 as builder

WORKDIR /work

COPY go.* *.go ./

RUN go mod download && CGO_ENABLED=0 GOOS=linux go build

FROM gcr.io/distroless/static-debian11
WORKDIR /
COPY --from=builder /work/tcp-connect-go /usr/bin/tcp-connect-go
