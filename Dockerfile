FROM golang:1.21-bullseye
LABEL authors="harshabose"

# Install build dependencies
RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    make \
    pkg-config \
    yasm \
    autoconf \
    automake \
    libtool \
    libass-dev \
    texinfo \
    nasm

# Create app directory
WORKDIR /app

# Clone the repository (you'll need to provide GITHUB_TOKEN at build time)
ARG GITHUB_TOKEN
RUN git clone --recursive https://${GITHUB_TOKEN}@github.com/harshabose/simple_skyline_sonata.git .

# Run the make commands to install dependencies
RUN make install-libx264
RUN make install-libopus
RUN make install-ffmpeg
RUN make install-mavp2p

# Build your application
RUN make build-delivery-gcs

# Run your application
CMD make run-delivery-gcs