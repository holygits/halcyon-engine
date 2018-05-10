FROM alpine:edge

RUN apk add --no-cache R \
                       R-dev \
                       build-base \
                       && R -e 'install.packages(c("packrat", "ggplot2", "devtools", "plyr"))'
