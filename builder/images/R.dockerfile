FROM alpine:edge

RUN apk add --no-cache R \
                       R-dev \
                       build-base \
    && R -e 'install.packages("packrat")'
