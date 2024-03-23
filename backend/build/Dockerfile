ARG GO_VERSION=1.22

# Builder
FROM golang:${GO_VERSION}-alpine AS build
RUN apk add --no-cache upx
RUN addgroup -S gofit && \
    adduser -S -u 1001 -g gofit gofit
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN go build \
    -installsuffix 'static' \
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
