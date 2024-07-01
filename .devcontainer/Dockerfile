FROM bamaas/devcontainer:2f03e7ccfb73871b7deb353e63081b806e79de94
SHELL ["/bin/zsh", "-c"]
USER root

# Setup
COPY . .
WORKDIR /setup

# Install node version
RUN source ~/.zshrc && \
    . $NVM_DIR/nvm.sh && \
    nvm install $(cat .nvmrc)

# Cleanup
WORKDIR /src
RUN rm -rf /setup

