ARG GO_VERSION=1.22

# Build the static files
FROM node:20.10.0-alpine3.19 as frontend-builder
WORKDIR /builder
COPY ./frontend/package.json /frontend/package-lock.json ./
RUN npm ci
COPY ./frontend .
RUN rm ./.env.production &&  \
    echo 'PUBLIC_BACKEND_BASE_URL="/api"' > .env.production
RUN npm run build

# Build the binary
FROM golang:${GO_VERSION}-alpine AS backend-builder
RUN apk add --no-cache upx
RUN addgroup -S gofit && \
    adduser -S -u 1001 -g gofit gofit
WORKDIR /src
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download
COPY ./ ./
COPY --from=frontend-builder /builder/build ./backend/internal/assets/static/
WORKDIR /src/backend
RUN go build \
    -ldflags="-s -w" \
    -o /app/gofit ./cmd/gofit
RUN upx --best --lzma /app/gofit
RUN mkdir /data

# Final
FROM scratch AS final
COPY --from=backend-builder /etc/passwd /etc/passwd
COPY --from=backend-builder /etc/group /etc/group
USER gofit
COPY --from=backend-builder --chown=gofit:gofit /data /data
COPY --from=backend-builder --chown=gofit:gofit /app/gofit /usr/local/bin/gofit
ENTRYPOINT ["/usr/local/bin/gofit"]
