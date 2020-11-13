FROM golang:1.14-alpine3.11 AS builder

ARG GO_LDFLAGS

COPY . /go/src/edgebus

RUN CGO_ENABLED=0 go build -v -o /usr/local/bin/server -ldflags "$GO_LDFLAGS -w -s" \
edgebus/server/cmd/server.go


FROM alpine:3.11

COPY --from=builder /usr/local/bin/server /usr/local/bin/server

ENTRYPOINT ["server"]