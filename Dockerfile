FROM golang:1.11

ENV GO111MODULE on

## Set up direnv
RUN apt-get update && apt-get install -y direnv 
COPY .envrc ./.envrc
RUN direnv allow

## Build
RUN mkdir -p ${GOPATH}/src/github.com/rerost/pubsub-duplicate-sample
RUN ln -s /app ${GOPATH}/src/github.com/rerost/pubsub-duplicate-sample
COPY . /app
WORKDIR /app
RUN make build
