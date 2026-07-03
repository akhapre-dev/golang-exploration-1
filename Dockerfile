# ─────────────────────────────────────────────────────────────────
# Multi-stage Dockerfile for golang-adk-exploration-1
# Stage 1: build a statically linked binary
# Stage 2: minimal distroless runtime image
# ─────────────────────────────────────────────────────────────────

# ── Stage 1: Build ───────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

# Install git so go mod download can fetch VCS metadata
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

# Cache dependency downloads as a separate layer
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source and compile a statically linked binary
COPY . .
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
      -ldflags "-s -w \
        -X main.version=${VERSION} \
        -X main.commit=${COMMIT} \
        -X main.buildDate=${BUILD_DATE}" \
      -o /out/agent \
      ./cmd/agent

# ── Stage 2: Runtime ─────────────────────────────────────────────
# gcr.io/distroless/static contains only CA certs + timezone data;
# no shell, no package manager — minimal attack surface.
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /out/agent /agent

# The ADK web UI listens on 8080 by default
EXPOSE 8080

# GOOGLE_API_KEY (and other config) are injected at runtime via --env-file or -e;
# they are intentionally NOT baked into this image.
ENTRYPOINT ["/agent"]
