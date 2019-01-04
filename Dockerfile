FROM golang:1.11

ENV GO111MODULE on

## Set up direnv
RUN apt-get update & apt-get install -y direnv 
COPY .envrc .evnrc
RUN direnv allow

## Build
RUN make build
