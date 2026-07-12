# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Install git for fetching private modules if needed.
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Ensure dashboard-dist exists for Go embed. CI provides the real build; local
# Docker builds get a placeholder so the embed directive does not fail.
RUN mkdir -p internal/ui/dashboard-dist && \
    (test -f internal/ui/dashboard-dist/index.html || \
     printf '%s\n' '<html><head><title>OpenMuara</title></head><body><p>Dashboard not built. Run <code>task ui:build</code> first.</p></body></html>' > internal/ui/dashboard-dist/index.html)

# Build a static binary so the final image can be minimal.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/muara ./cmd/muara

# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

# Create a non-root user and ensure /app is writable for config/data.
RUN addgroup -g 1000 muara && \
    adduser -u 1000 -G muara -D -s /sbin/nologin muara && \
    mkdir -p /app/.muara && \
    chown -R muara:muara /app

WORKDIR /app

COPY --from=builder /bin/muara /usr/local/bin/muara
COPY --from=builder /src/scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
COPY --from=builder /src/plugins /app/plugins
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

EXPOSE 9000

USER muara:muara

HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
  CMD ["muara", "health"]

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["start"]
