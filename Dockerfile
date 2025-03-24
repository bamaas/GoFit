# Arguments (passed automatically by Makefile via 'make image/build')
ARG GO_VERSION
ARG NODE_VERSION

# Build the static files
FROM node:${NODE_VERSION}-alpine as frontend-builder
WORKDIR /src
COPY ./frontend/package.json /frontend/package-lock.json ./
RUN npm ci
COPY ./frontend .
RUN echo 'PUBLIC_BACKEND_BASE_URL="/api"' > .env.production
RUN npm run build

# Build the binary
FROM golang:${GO_VERSION}-alpine AS backend-builder
RUN apk add --no-cache upx && \
    addgroup -S gofit && \
    adduser -S -u 1001 -g gofit gofit
WORKDIR /src
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download
COPY ./backend .
COPY --from=frontend-builder /src/build ./internal/assets/static/
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o /bin/gofit ./cmd/gofit && \
    upx --best --lzma /bin/gofit && \
    mkdir -p /data

# Final
FROM scratch AS final
ARG GOFIT_VERSION

LABEL org.opencontainers.image.source="https://github.com/bamaas/GoFit"
LABEL org.opencontainers.image.description="A weight tracking app"
LABEL org.opencontainers.image.version="${GOFIT_VERSION}"

COPY --from=backend-builder /etc/passwd /etc/passwd
COPY --from=backend-builder /etc/group /etc/group
COPY --from=backend-builder --chown=gofit:gofit /data /data
COPY --from=backend-builder --chown=gofit:gofit /bin/gofit /usr/local/bin/gofit

# Run
USER gofit
ENTRYPOINT ["/usr/local/bin/gofit"]
