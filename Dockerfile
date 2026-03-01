# syntax=docker/dockerfile:1

# RunNotes Docker Desktop Extension
# Multi-stage build: backend + frontend

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

# TODO: Add backend build stage
# TODO: Add frontend build stage
# TODO: Copy ui assets

COPY metadata.json .
