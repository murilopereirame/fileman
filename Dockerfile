FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /fileman

FROM build AS tests
RUN go test -v ./...

FROM alpine:latest

ARG PUID=1000
ARG PGID=1000

WORKDIR /app
COPY --from=build /fileman ./

RUN addgroup -g "${PGID}" -S go \
          && adduser -u "${PUID}" -G go -S -D -h /home/go go

USER go:go

ENTRYPOINT ["/app/fileman"]