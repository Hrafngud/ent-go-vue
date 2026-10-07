# syntax=docker/dockerfile:1
FROM golang:1.27.1-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM arigaio/atlas:1.3.0 AS migrations
WORKDIR /app
COPY atlas.hcl ./
COPY ent/migrate/migrations ./ent/migrate/migrations
USER 10001:10001
ENTRYPOINT ["/atlas"]
CMD ["migrate", "apply", "--env", "app"]

FROM alpine:3.23 AS runtime
RUN apk add --no-cache ca-certificates \
    && addgroup -S -g 10001 app \
    && adduser -S -D -H -u 10001 -G app app
WORKDIR /app
COPY --from=builder /out/api ./api
USER app
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/api"]
