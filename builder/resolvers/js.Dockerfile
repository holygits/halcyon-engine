FROM alpine:edge

RUN apk add --no-cache yarn \
   # Pre-emptivley install some popular packages
    && yarn install lodash \
                    request \
                    moment \
                    d3
