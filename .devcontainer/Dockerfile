ARG DEVCONTAINER_TAG=latest-minimal
FROM bamaas/devcontainer:${DEVCONTAINER_TAG}
COPY ./mise.toml mise.toml
RUN mise trust && \
    mise install
# WORKDIR /builder
# COPY  . .
# RUN make development/setup
# WORKDIR /src
# RUN rm -rf /builder