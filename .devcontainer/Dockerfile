FROM golang:alpine

ARG USER=vscode
ARG UID=1000
ARG GID=1000

RUN apk add -q --update sudo openssh-client git zsh starship

RUN adduser $USER -s /bin/zsh -D -u $UID $GID && \
    mkdir -p /etc/sudoers.d && \
    echo $USER ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USER && \
    chmod 0440 /etc/sudoers.d/$USER

USER $USER

RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest

RUN echo "eval \"$(starship init zsh)\"" >> /home/$USER/.zshrc

USER root