FROM golang:1.18-alpine3.15 as builder

WORKDIR /app

ADD . .

RUN go install ./cmd/server

FROM alpine:3.15

COPY --from=builder /go/bin/server /server


ENTRYPOINT /server
