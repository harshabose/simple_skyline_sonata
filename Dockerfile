FROM golang:1.23-bullseye
LABEL authors="harshabose"

SHELL ["/bin/bash", "-c"]

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
WORKDIR /simple_skyline_sonata-gcs

COPY . .

# Run the make commands to install dependencies
RUN make install-ffmpeg-docker
RUN make install-mavp2p

# Will create complete.env file
RUN make create-env-file

RUN bash -c 'make build-delivery-gcs'

# Run main application
CMD bash -c 'make run-delivery-gcs'
