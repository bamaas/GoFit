FROM jdxcode/mise:2025.3.6
WORKDIR /src
USER root
COPY ./mise.toml mise.toml
COPY ./Makefile Makefile
RUN make dev/install/tools && \
    rm -rf Makefile

# Install Docker CLI
RUN apt-get update && \
    apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release && \
    install -m 0755 -d /etc/apt/keyrings && \
    curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg && \
    chmod a+r /etc/apt/keyrings/docker.gpg && \
    echo \
        "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
        "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
        tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt-get update && \
    apt-get install -y docker-ce-cli --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Install bsdmainutils
RUN apt-get update && \
    apt-get install -y bsdmainutils --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

ENTRYPOINT [ "" ]