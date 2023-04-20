FROM golang:1.20-alpine3.17 as builder

ENV GOMODCACHE=/cache/gomod
ENV GOCACHE=/cache/gobuild

WORKDIR /app

ADD . .

RUN --mount=type=cache,target=/cache/gomod \
    go mod download

RUN --mount=type=cache,target=/cache/gomod \
    --mount=type=cache,target=/cache/gobuild,sharing=locked \
    go install ./cmd/server

FROM alpine:3.17

COPY --from=builder /go/bin/server /server


ENTRYPOINT /server
