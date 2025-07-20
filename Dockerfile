# syntax=docker/dockerfile:1
FROM ubuntu:24.04

SHELL ["/bin/bash", "-c"]

RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get -y install git
RUN apt-get -y install git
RUN apt-get -y install make
RUN apt-get -y install wget
RUN apt-get -y install curl
RUN apt-get -y install nano

# Install Go
RUN wget https://go.dev/dl/go1.24.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.24.5.linux-amd64.tar.gz
RUN rm go1.24.5.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
ENV PATH="$HOME/go/bin:${PATH}"

RUN go install -v golang.org/x/tools/gopls@latest
RUN go install -v github.com/air-verse/air@latest

# Install Node & NPM
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash
RUN \. "$HOME/.nvm/nvm.sh" && nvm install 22