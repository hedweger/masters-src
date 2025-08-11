FROM debian:bookworm as builder

RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    wget \
    curl \
    pkg-config \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN git clone --branch masters --depth 1 https://github.com/hedweger/libiec61850.git
WORKDIR /libiec61850/third_party/mbedtls
run wget https://github.com/Mbed-TLS/mbedtls/archive/refs/tags/v2.28.8.tar.gz && tar xzf v2.28.8.tar.gz && mv mbedtls-2.28.8 mbedtls-2.28
workdir /libiec61850
RUN make INSTALL_PREFIX=/usr/local WITH_MBEDTLS=1 install
RUN rm -rf /libiec61850 && ldconfig


WORKDIR /app
COPY . .
RUN mkdir -p build && cd build && cmake .. && make
CMD ["./build/bin/ied-client"]
