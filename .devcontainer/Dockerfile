FROM bamaas/devcontainer:b7c3ac34f282e0f95197d5eaaa16ff715380efa5
COPY ./mise.toml mise.toml
RUN /bin/zsh -c 'set -o pipefail && eval "$(~/.local/bin/mise activate zsh)" && mise trust && mise install'