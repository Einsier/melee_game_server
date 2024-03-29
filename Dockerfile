FROM golang:1.17-alpine as builder
WORKDIR /root/go/src/github.com/einsier/ustc_melee_game
COPY . /root/go/src/github.com/einsier/ustc_melee_game
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
RUN go build -o game-server run.go

FROM alpine:latest
# environment variable
ARG ETCD_ADDR
ENV ENV_ETCD_ADDR=${ETCD_ADDR} \
    ENV_PLAYER_NUM=10
WORKDIR  /root/go/src/github.com/einsier/ustc_melee_game
COPY --from=builder  /root/go/src/github.com/einsier/ustc_melee_game/game-server .
EXPOSE 8000/tcp
EXPOSE 8001/tcp
EXPOSE 8080/tcp
ENTRYPOINT ./game-server -etcdAddr ${ENV_ETCD_ADDR} -playerNum ${ENV_PLAYER_NUM}