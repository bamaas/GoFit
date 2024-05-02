ARG GO_VERSION=1.22

# Stage 1: Build the static files
FROM node:20.10.0-alpine3.19 as frontend-builder
WORKDIR /builder
COPY ./frontend/package.json /frontend/package-lock.json ./
RUN npm ci
COPY ./frontend .
RUN rm ./.env.production &&  \
    echo 'PUBLIC_BACKEND_BASE_URL=""' > .env.production
RUN npm run build

# Builder
FROM golang:${GO_VERSION}-alpine AS build
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
    -o /app ./cmd/gofit
RUN upx --best --lzma /app

# Final
FROM scratch AS final
WORKDIR /app
COPY --from=build /app /gofit
COPY --from=build /etc/passwd /etc/passwd
USER gofit
ENTRYPOINT ["/gofit"]
