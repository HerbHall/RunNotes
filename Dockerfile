# syntax=docker/dockerfile:1

# RunNotes Docker Desktop Extension
# Multi-stage build: frontend + backend

# Stage 1: Build React UI
FROM --platform=$BUILDPLATFORM node:22-alpine AS ui-builder
WORKDIR /app
COPY ui/package.json ui/package-lock.json* ./
RUN npm ci
COPY ui/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
ARG TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o bin/backend ./cmd/backend

# Stage 3: Final extension image
FROM alpine:3.21
LABEL org.opencontainers.image.title="RunNotes" \
      org.opencontainers.image.description="Attach notes and annotations to your Docker containers" \
      org.opencontainers.image.vendor="Herb Hall" \
      com.docker.desktop.extension.api.version=">= 0.3.3" \
      com.docker.desktop.extension.icon="https://raw.githubusercontent.com/HerbHall/RunNotes/main/docker.svg" \
      com.docker.extension.screenshots="" \
      com.docker.extension.detailed-description="RunNotes lets you attach notes to containers so you never forget why a container exists or what you were testing." \
      com.docker.extension.publisher-url="https://github.com/HerbHall" \
      com.docker.extension.changelog=""

COPY --from=builder /app/bin/backend /backend
COPY docker-compose.yaml .
COPY metadata.json .
COPY docker.svg .
COPY --from=ui-builder /app/build ui

CMD ["/backend", "-socket", "/run/guest-services/backend.sock"]
