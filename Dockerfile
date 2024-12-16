# Setup the build environment
ARG DEVCONTAINER_TAG=placeholder_parsed_from_devcontainer_dockerfile
FROM bamaas/devcontainer:${DEVCONTAINER_TAG} AS builder
WORKDIR /builder
COPY . .
RUN make development/setup

# Build the binary
RUN make binary/build binary/compress
RUN mkdir /data

# Create run user
FROM alpine:3.21 AS user
RUN apk add --no-cache upx
RUN addgroup -S gofit && \
    adduser -S -u 1001 -g gofit gofit

# Final
FROM scratch AS final
COPY --from=user /etc/passwd /etc/passwd
COPY --from=user /etc/group /etc/group
USER gofit
COPY --from=builder --chown=gofit:gofit /data /data
COPY --from=builder --chown=gofit:gofit /builder/backend/bin/gofit /usr/local/bin/gofit
ENTRYPOINT ["/usr/local/bin/gofit"]
