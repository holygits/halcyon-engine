FROM alpine:edge

RUN apk add --no-cache python3 \
                       python3-dev \
                       build-base \
                       python3-pip \
                       git \
    && pip install setuptools wheel

# pandas is a very common python package but there is no wheel for alpine, hence we pre-emptivley build pandas and include it in the initial cache
RUN git clone --depth=1 https://github.com/pandas-dev/pandas \
    && cd pandas \
    && pip setup.py bdist_wheel \
    && pip install dist/

# TODO: Preemptivley install other popular python packages
RUN pip install requests \
                pymongo \
                scipy \
                lxml \
                psycopg2
