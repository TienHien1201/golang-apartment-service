# syntax=docker/dockerfile:1.7

# ── Stage 1: Builder ─────────────────────────────────────────
# Dùng BuildKit cache mounts: module cache + build cache
# → lần build đầu ~2-3 phút, lần sau (code thay đổi) ~15-30 giây
FROM golang:1.24-alpine AS builder

ARG APP_VERSION=dev
ARG GIT_COMMIT=unknown

WORKDIR /app

# Copy manifests riêng → layer này cache đến khi go.mod/go.sum thay đổi
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

COPY . .

# -trimpath : loại bỏ abs path ra khỏi binary (security + nhỏ hơn ~5 MB)
# -s -w     : strip debug info + DWARF (~30-40% nhỏ hơn)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
        -trimpath \
        -ldflags="-s -w \
            -X main.Version=${APP_VERSION} \
            -X main.GitCommit=${GIT_COMMIT}" \
        -o /out/server \
        ./cmd/app

# ── Stage 2: Runtime ─────────────────────────────────────────
# alpine:3.21 (~8 MB) + tzdata + ca-certs (~12 MB) + binary stripped (~20-28 MB)
# Tổng: ~35-45 MB  vs  ~70-80 MB trước đây
FROM alpine:3.21 AS runtime

# Non-root user — không chạy process bằng root trong container
RUN addgroup -S appgrp && adduser -S appuser -G appgrp

# tzdata        : cần cho mysql DSN loc=Local + TZ=Asia/Ho_Chi_Minh
# ca-certificates: cần cho HTTPS (Cloudinary, OAuth, SMTP, ES)
RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=builder --chown=appuser:appgrp /out/server ./server
# Config được copy làm fallback; docker-compose volume mount sẽ override
COPY --from=builder --chown=appuser:appgrp /app/config ./config

USER appuser

EXPOSE 1424

HEALTHCHECK --interval=15s --timeout=5s --start-period=30s --retries=3 \
    CMD wget -qO- http://localhost:1424/api/health || exit 1

ENTRYPOINT ["./server"]
# Override bằng command: ["--env=prod"] trong docker-compose
CMD ["--env=dev"]
