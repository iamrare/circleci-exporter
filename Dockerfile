FROM golang:1.14.7-alpine3.11 AS builder
LABEL maintainer="iamrare"

ENV GO111MODULE=on

COPY ./ /go/src/github.com/iamrare/circleci-exporter
WORKDIR /go/src/github.com/iamrare/circleci-exporter

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -o /bin/main

FROM alpine:3.12.0

RUN apk --no-cache add ca-certificates \
     && addgroup exporter \
     && adduser -S -G exporter exporter

USER exporter

COPY --from=builder /bin/main /bin/main
ENV URL "https://circleci.com/api/v2/insights/gh/yourOrg/yourRepo/workflows/deploy"
ENV AUTH_TOKEN "yourToken"
ENV LISTEN_PORT=9179
EXPOSE 9179
ENTRYPOINT ["/bin/main"]
