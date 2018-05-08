FROM alpine:edge

RUN apk add --no-cache python3 \
                       python3-dev \
                       build-base \
                       python3-pip
