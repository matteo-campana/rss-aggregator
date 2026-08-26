# syntax=docker/dockerfile:1

# --- build -------------------------------------------------------------------
FROM golang:1.24-alpine AS build

WORKDIR /src

# Dependencies first: they change far less often than the code, so this layer
# survives most rebuilds.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# --- migrations --------------------------------------------------------------
# A separate image so migrations run as their own one-shot step instead of at
# API startup, where two replicas would race each other.
FROM golang:1.24-alpine AS goose

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.21.1

FROM alpine:3.20 AS migrate

RUN apk add --no-cache ca-certificates

COPY --from=goose /go/bin/goose /usr/local/bin/goose
COPY sql/schema /migrations/schema

WORKDIR /migrations
USER nobody

# The connection comes from GOOSE_DRIVER/GOOSE_DBSTRING rather than positional
# arguments: goose reads the command from the wrong argument when both are given.
ENV GOOSE_DRIVER=postgres
ENTRYPOINT ["goose", "-dir", "schema"]
CMD ["up"]

# --- api ---------------------------------------------------------------------
FROM alpine:3.20 AS api

# wget comes from busybox and is what the healthcheck uses.
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 app

COPY --from=build /out/api /usr/local/bin/api

USER app
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --spider "http://127.0.0.1:${PORT:-8080}/api/v1/health" || exit 1

ENTRYPOINT ["api"]
