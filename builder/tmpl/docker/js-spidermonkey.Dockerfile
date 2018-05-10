FROM frolvlad/alpine-glibc

RUN wget https://raw.githubusercontent.com/orctom/alpine-glibc-packages/master/usr/lib/libstdc++.so.6.0.21 \
         -O /usr/lib/libstdc++.so.6.0.21 \
    && cp /usr/lib/libstdc++.so.6.0.21 /usr/lib/libstdc++.so.6

ADD lib/js /usr/bin/js
