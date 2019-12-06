ARG GO_VERSION=1.11
FROM golang:${GO_VERSION}-alpine AS builder
ARG DOCKER_GIT_CREDENTIALS
RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .
FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates
RUN apk add bash
RUN mkdir /app
COPY --from=builder /app /app
EXPOSE 8080
WORKDIR /app
CMD ["./app"]