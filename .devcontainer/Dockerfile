FROM bamaas/devcontainer:0c9e378410cf10210b512eb8666111f7c3cef519
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

