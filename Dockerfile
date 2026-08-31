# Build frontend
FROM node:24-alpine AS frontend-builder

WORKDIR /app/web

# Install dependencies
COPY web/package*.json ./
RUN npm ci

# Build
COPY web/ ./
RUN npm run build

# Build backend
FROM golang:1.27.0-alpine AS backend-builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
ARG VERSION=dev
COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w -X github.com/ArtemStepanov/caddy-admin-ui/internal/version.Version=${VERSION}" -o caddy-admin-ui ./cmd/server

# Final image
FROM alpine:3.24

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs tzdata \
    && addgroup -S caddy-admin-ui \
    && adduser -S -G caddy-admin-ui -h /app caddy-admin-ui

WORKDIR /app

# Copy binary
COPY --from=backend-builder --chown=caddy-admin-ui:caddy-admin-ui /app/caddy-admin-ui .

# Copy frontend
COPY --from=frontend-builder --chown=caddy-admin-ui:caddy-admin-ui /app/web/dist ./web/dist

# Create data directory
RUN mkdir -p /app/data && chown caddy-admin-ui:caddy-admin-ui /app/data

# Environment variables
ENV GIN_MODE=release
ENV DB_PATH=/app/data/routes.db
ENV CADDY_ADMIN_URL=http://localhost:2019
ENV LISTEN_ADDR=127.0.0.1:3000

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/healthz || exit 1

# Run
USER caddy-admin-ui
CMD ["./caddy-admin-ui"]
