# Build Geth in a stock Go builder container
FROM golang:1.16-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git bash

ADD . /3f-chain
RUN cd /3f-chain && make 3fnode

# Pull 3fnode into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates curl jq tini
COPY --from=builder /3f-chain/build/bin/3fnode /usr/local/bin/

EXPOSE 8545 8546 8547 30303 30303/udp
ENTRYPOINT ["3fnode"]