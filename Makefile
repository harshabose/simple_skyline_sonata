ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    NULL_DEV := NUL
    export PKG_CONFIG := pkgconfiglite
else
    DETECTED_OS := $(shell uname -s)
    NULL_DEV := /dev/null
endif

# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# COMMAND LINE ARGUMENTS
# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

# Default value for cgo_enabled
CGO_ENABLED ?= false

# Validate CGO_ENABLED argument
ifneq ($(CGO_ENABLED),true)
ifneq ($(CGO_ENABLED),false)
$(error CGO_ENABLED must be either 'true' or 'false'. Usage: make target CGO_ENABLED=true)
endif
endif

# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR-VAR
# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

# GENERAL VARIABLES
SHELL := /bin/bash
WORKING_DIR := $(shell pwd)
THIRD_PARTY_DIR := $(WORKING_DIR)/third_party

# GO RELATED VARIABLES
GOBASE=$(WORKING_DIR)
GOBIN=$(GOBASE)/

# VARIABLES FOR BUILD
CMD_DIR := $(WORKING_DIR)/cmd
BUILD_DIR := $(WORKING_DIR)/build
VERSION=$(shell git describe --tags --always --long --dirty)
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# VARIABLES FOR FFMPEG CHECK
FFMPEG_VERSION := n7.0
FFMPEG_DIRECTORY := $(THIRD_PARTY_DIR)/ffmpeg
FFMPEG_SRC_DIR := $(FFMPEG_DIRECTORY)/src

# VARIABLES FOR X264
X264_DIRECTORY := $(THIRD_PARTY_DIR)/x264
X264_SRC_DIR := $(X264_DIRECTORY)/src

# VARIABLES FOR VPX
VPX_DIRECTORY := $(THIRD_PARTY_DIR)/vpx
VPX_SRC_DIR := $(VPX_DIRECTORY)/src

OPUS_DIRECTORY := $(THIRD_PARTY_DIR)/libopus
OPUS_SRC_DIR := $(OPUS_DIRECTORY)/src

# VARIABLES FOR MAVP2P
MAVP2P_INSTALL_DIR := $(THIRD_PARTY_DIR)/mavp2p

# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV
# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

# COMPILE ENVIRONMENT (only when CGO is enabled)
ifeq ($(CGO_ENABLED),true)
COMPILE_ENV := CGO_LDFLAGS="-L$(FFMPEG_DIRECTORY)/lib -L$(X264_DIRECTORY)/lib -L$(VPX_DIRECTORY)/lib" \
               CGO_CFLAGS="-I$(FFMPEG_DIRECTORY)/include -I$(X264_DIRECTORY)/include -I$(VPX_DIRECTORY)/include" \
               PKG_CONFIG_PATH="$(FFMPEG_DIRECTORY)/lib/pkgconfig:$(X264_DIRECTORY)/lib/pkgconfig:$(VPX_DIRECTORY)/lib/pkgconfig" \
               LD_LIBRARY_PATH="$(FFMPEG_DIRECTORY)/lib:$(X264_DIRECTORY)/lib:$(VPX_DIRECTORY)/lib:$$LD_LIBRARY_PATH"

ifeq ($(DETECTED_OS),Darwin)
	CGO_RUNTIME_ENV_MACOS := DYLD_LIBRARY_PATH="$(FFMPEG_DIRECTORY)/lib:$(X264_DIRECTORY)/lib:$(VPX_DIRECTORY)/lib"
endif

else
COMPILE_ENV :=
CGO_RUNTIME_ENV_MACOS :=
endif

# BASE ENVIRONMENT VARIABLES
FIREBASE_ENV := FIREBASE_TYPE=service_account \
                FIREBASE_PROJECT_ID=iitb-rgstc-signalling-server \
                FIREBASE_CLIENT_EMAIL=firebase-adminsdk-s07hu@iitb-rgstc-signalling-server.iam.gserviceaccount.com \
                FIREBASE_PRIVATE_KEY_ID=a512e26c961557a4d97498ea9b00d84ce683dce8 \
                FIREBASE_PRIVATE_KEY='-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDB0kyQDvmEUj5K\nJsfTGX5DYCFpULBUr0kuNyTVFzeRXDuTWKKOotk9qo8VEmCpQcFviiayk9piUCWW\nfSONwoXnEP+GI1SCl9N2zMzvuuuWPWH9xgdHRpdpWEHsrtL6DuoVpepE99uQ/yt4\nj+9QUCDEryyQ4MUmE8DjerQEM+Tj4VgozW99dZ5MzjgLTwrJhyHeljliMrG+SB1a\nRfUNsxzk/UrthSGIgvGvE2Fjk9DVtIWWbgBfVDTB+DdvYJFn3Xg10C4nmWxGjyh1\nGt2AyIW9asl8oOKljUF5LzJUaniFVvRiqgZolBkCxjDHjjnENJKqiqOTxbVAzibL\npSnjwzx3AgMBAAECggEAIPv5/ZYezm70nMfmv70Z6LtmVDbgGzlNWekWgpEV6s3o\ncZXm7CE4mS76dJqRCpzfH21CUqeoxYxgKTEYqNpO0VjqM1i13BecbB5EThPgXcwK\nbhaSTIXt5IaZiX7i9p0tJwv6R0xq+E0Eh9ru3hsUyIQLMIif5G/+JnhORFzUehcm\niWfgPtbh6RMiRoarrI4p+cFFqEZheJVrqsI+StPYpVCYZi4YLArhryjNVMsXhUsV\nhbXDDTIKPaLWsU/Ct0+vjF8CUBWaJCZ5z9lqLd0gnQn3XzVHqTW9jQ+o5LYkE3x8\n6Uw+JXBaUt0M09c5/eF3QWI7Jg+YMKwSVjHFPCpyuQKBgQDf5HVfRLwlZqm9FPK4\nmXawWvgYc7xq5s8Rt/wk9Z2I/y8E6Gw5uJjgw//NnCvET0y0i9JpC/kB8Nu6Xj4V\nzx1j6xFRPaPwRmQH2b5r3otoWuBhwffaehSw+aSovEaWSUXlmbN/eJKpkzHU5ay1\nVQB72UUYrRvntlTYmwQeHOMnXQKBgQDdneBBGJYy0rMsicTelsMCWjQ2PiHzlx0i\nTs8CulDyjdew5XqiH08i75Buz2hRMWUcvyc49wMpdZIvtdae92DTWvX7QZjiihrJ\nDg3rX15Fhtb38RdHyRssGrOl7u0q5BKzhuY+Lq4YwjQTgzg3Zmy8M10Uo1i/mh9C\nP4Ae8chZ4wKBgCKxOMqxUOIOvWByHYYjKXP8NJM9Y8XAy/c35hcoA+gVeoitJw/u\nnam+VSXb/CAoFX+oZssmMshtNO706XPhqvEvnHhVL9DsZ1WcFNiMHFfoNPqQ3sH4\nxroBhNUsj1d8NRt1rI2k9jzWdRNDH3bdm/yU1xMSx88ovo7tvj6YRU51AoGBAMkT\nwPBvZYBRin6DthucQO32eF8q+tUwrB9/z/YSpPWe2zBG1oEY1U3GfY79IxJgNfTi\nP61A+h547Z3aaBQuMi0y3/MMLrKFSg5YcSq5iiidUpj+p/fbMYtP4uZQpeH/tDQt\n1uReqFoQgv2dVrl1dn1AQVlDaHfYWDpcsVviVr2vAoGBAJyURZif4b2nRU82uXXn\n+IdIDXSAyta640GKjwbYRydq295mVC1mCRFYpTq7D61XDGhbSLO+cb0mA4CGzNrN\nQebL2yXn83gGpQiiJ/dy3uIMLnk22iWTF7GTfH7sjkRDMrsUqdezQ3kQepaeLi6Y\nCqzfYevkBBh8joHXIHsC7BDI\n-----END PRIVATE KEY-----\n' \
                FIREBASE_CLIENT_ID=106924326990810690130 \
                FIREBASE_AUTH_URI=https://accounts.google.com/o/oauth2/auth \
                FIREBASE_AUTH_TOKEN_URI=https://oauth2.googleapis.com/token \
                FIREBASE_AUTH_PROVIDER_X509_CERT_URL=https://www.googleapis.com/oauth2/v1/certs \
                FIREBASE_AUTH_CLIENT_X509_CERT_URL=https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-s07hu%40iitb-rgstc-signalling-server.iam.gserviceaccount.com \
                FIREBASE_UNIVERSE_DOMAIN=googleapis.com

NETWORK_ENV := STUN_SERVER_URL=stun:stun.skyline-sonata.in:3478 \
               TURN_UDP_SERVER_URL=turn:turn.skyline-sonata.in:3478 \
               TURN_TCP_SERVER_URL=turn:turn.skyline-sonata.in:3478 \
               TURN_TLS_SERVER_URL=turn:turn.skyline-sonata.in:5349 \
               TURN_SERVER_USERNAME=super@skyline-sonata \
               TURN_SERVER_PASSWORD=rufryz-wofdI5-mawged

HARDWARE_ENV := MAVP2P_EXE_PATH=$(MAVP2P_INSTALL_DIR)/mavp2p \
                MAVLINK_SERIAL=/dev/ttyTHS0:115200

# Runtime environments
CGO_RUNTIME_ENV := $(HARDWARE_ENV) $(FIREBASE_ENV) $(NETWORK_ENV) $(COMPILE_ENV) $(CGO_RUNTIME_ENV_MACOS)
RUNTIME_ENV := $(HARDWARE_ENV) $(FIREBASE_ENV) $(NETWORK_ENV)


# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# ENVIRONMENT FILE CREATION
# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

create-env-file:
	@echo "Creating environment files (without CGO compilation flags)..."
	@# Create runtime environment file
	@echo "# Runtime Environment Variables" > runtime.env
	@echo "$(HARDWARE_ENV)" | sed 's/ \([A-Z_]*=\)/;\1/g' >> runtime.env
	@echo "" >> runtime.env
	@echo "# Firebase Configuration" >> runtime.env
	@echo "$(FIREBASE_ENV)" | sed 's/ \([A-Z_]*=\)/;\1/g' >> runtime.env
	@echo "" >> runtime.env
	@echo "# Network Configuration" >> runtime.env
	@echo "$(NETWORK_ENV)" | sed 's/ \([A-Z_]*=\)/;\1/g' >> runtime.env
	@echo "Runtime environment file created at runtime.env (without CGO compilation flags)"

create-env-file-with-cgo-enabled:
ifeq ($(CGO_ENABLED),false)
	@echo "ERROR: CGO_ENABLED must be set to 'true' to create environment file with CGO flags"
	@echo "Usage: make create-env-file-with-cgo-enabled CGO_ENABLED=true"
	@exit 1
endif
	@echo "Creating environment files with CGO compilation flags..."
	@# Create runtime environment file
	@echo "# Runtime Environment Variables" > runtime.env
	@echo "$(HARDWARE_ENV)" | sed 's/ \([A-Z_]*=\)/;\1/g' >> runtime.env
	@echo "" >> runtime.env
	@echo "# Firebase Configuration" >> runtime.env
	@echo "$(FIREBASE_ENV)" | sed 's/ \([A-Z_]*=\)/;\1/g' >> runtime.env
	@echo "" >> runtime.env
	@echo "# Network Configuration" >> runtime.env
	@echo "$(NETWORK_ENV)" | sed 's/ \([A-Z_]*=\)/;\1/g' >> runtime.env
	@echo "" >> runtime.env
	@echo "# CGO Compilation Environment" >> runtime.env
	@echo "$(COMPILE_ENV)" | sed 's/ \([A-Z_]*=\)/;\1/g' >> runtime.env
	@echo "Runtime environment file created at runtime.env (with CGO compilation flags)"
	@# Create compilation environment file
	@echo "# Compilation Environment Variables" > compile.env
	@echo "$(COMPILE_ENV)" | sed 's/ \([A-Z_]*=\)/;\1/g' >> compile.env
	@echo "Compilation environment file created at compile.env"

# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# CGO DEPENDENCY VALIDATION
# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

debug-os:
	@echo "Detected OS: $(DETECTED_OS)"
	@echo "CGO Enabled: $(CGO_ENABLED)"
ifeq ($(DETECTED_OS),Darwin)
	@echo "Darwin-specific paths will be set"
else
	@echo "Non-Darwin paths will be set"
endif

check-cgo-enabled:
ifeq ($(CGO_ENABLED),false)
	@echo "ERROR: This target requires CGO_ENABLED=true"
	@echo "Usage: make $@ CGO_ENABLED=true"
	@exit 1
endif

# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

.PHONY: check help

help:
	@echo "Available targets:"
	@echo ""
	@echo "Basic targets:"
	@echo "  check                           - Check if required tools are installed"
	@echo "  create-env-file                 - Create environment files without CGO flags"
	@echo "  create-env-file-with-cgo-enabled - Create environment files with CGO flags (requires CGO_ENABLED=true)"
	@echo ""
	@echo "CGO-dependent targets (require CGO_ENABLED=true):"
	@echo "  install-third-party             - Install FFmpeg and related dependencies"
	@echo "  install-libx264                 - Install libx264"
	@echo "  install-libvpx-darwin           - Install libvpx for macOS"
	@echo "  install-libopus                 - Install libopus"
	@echo "  install-ffmpeg-linux            - Install FFmpeg for Linux"
	@echo "  install-ffmpeg-darwin           - Install FFmpeg for macOS"
	@echo "  install-ffmpeg-docker           - Install FFmpeg for Docker"
	@echo "  build-fpv-drone            - Build fpv drone (requires CGO)"
	@echo "  run-fpv-drone              - Run fpv drone (requires CGO)"
	@echo ""
	@echo "Non-CGO targets:"
	@echo "  install-windows-deps            - Install Windows dependencies"
	@echo "  install-mavp2p                  - Install mavp2p"
	@echo "  build-fpv-gcs              - Build fpv GCS"
	@echo "  run-fpv-gcs                - Run fpv GCS"
	@echo ""
	@echo "Usage with CGO:"
	@echo "  make install-third-party CGO_ENABLED=true"
	@echo "  make build-fpv-drone CGO_ENABLED=true"
	@echo ""
	@echo "Usage without CGO:"
	@echo "  make build-fpv-gcs"
	@echo "  make install-mavp2p"

check:
	git --version >$(NULL_DEV) 2>&1 || (echo "git is not installed or not in PATH"; exit 1)
	go version >$(NULL_DEV) 2>&1 || (echo "go is not installed or not in PATH"; exit 1)

# CGO-dependent installation targets
install-third-party: check-cgo-enabled install-ffmpeg-linux install-mavp2p

install-libx264: check-cgo-enabled
	mkdir -p $(X264_SRC_DIR)
	cd $(X264_SRC_DIR) && git clone https://code.videolan.org/videolan/x264.git .
	cd $(X264_SRC_DIR) && git checkout stable
	cd $(X264_SRC_DIR) && ./configure \
            --prefix=$(X264_DIRECTORY) \
            --enable-shared \
            --enable-pic
	cd $(X264_SRC_DIR) && make -j$(nproc)
	cd $(X264_SRC_DIR) && make install

install-libvpx-darwin: check-cgo-enabled
	echo "Installing libvpx for macOS ARM64..."
	mkdir -p $(VPX_DIRECTORY)
	mkdir -p $(VPX_SRC_DIR)
	cd $(VPX_DIRECTORY) && rm -rf $(VPX_SRC_DIR)
	mkdir -p $(VPX_SRC_DIR)
	echo "Cloning libvpx (this may take several minutes)..."
	cd $(VPX_SRC_DIR) && git clone https://chromium.googlesource.com/webm/libvpx .
	cd $(VPX_SRC_DIR) && git checkout v1.14.0
	cd $(VPX_SRC_DIR) && MACOSX_DEPLOYMENT_TARGET=10.15 ./configure \
		--target=arm64-darwin20-gcc \
		--prefix=$(VPX_DIRECTORY) \
        --enable-shared \
        --disable-static \
        --enable-pic \
        --enable-vp8 \
        --enable-vp9 \
      	--enable-vp8-encoder \
        --enable-vp9-encoder \
        --enable-vp8-decoder \
        --enable-vp9-decoder \
        --enable-runtime-cpu-detect \
        --enable-multithread \
        --disable-examples \
        --disable-tools \
        --disable-docs \
        --disable-unit-tests \
        --disable-debug \
        --enable-optimizations \
        --extra-cflags="-mmacosx-version-min=10.15" \
        --extra-cxxflags="-mmacosx-version-min=10.15"
		cd $(VPX_SRC_DIR) && make clean
		cd $(VPX_SRC_DIR) && make -j$(shell sysctl -n hw.ncpu)
		cd $(VPX_SRC_DIR) && make install
	if [ ! -d "$(VPX_DIRECTORY)/lib" ]; then \
		echo "libvpx installation failed: lib directory not found"; \
		exit 1; \
	fi
	echo "libvpx installation complete."

install-libopus: check-cgo-enabled
	mkdir -p $(OPUS_SRC_DIR)
	cd $(OPUS_SRC_DIR) && git clone https://github.com/xiph/opus.git .
	cd $(OPUS_SRC_DIR) && ./autogen.sh
	cd $(OPUS_SRC_DIR) && ./configure \
		--prefix=$(OPUS_DIRECTORY) \
        --enable-shared \
        --enable-pic \
        --enable-custom-modes \
        CFLAGS="-march=native -O3"
	cd $(OPUS_SRC_DIR) && make -j$(nproc)
	cd $(OPUS_SRC_DIR) && make install

# Non-CGO installation targets
install-windows-deps:
	if [ "$(DETECTED_OS)" = "Windows" ]; then \
		echo "Installing Windows dependencies..."; \
		pacman -Syu --noconfirm; \
		pacman -S --noconfirm --needed git diffutils mingw-w64-x86_64-toolchain pkg-config make yasm; \
	fi

install-ffmpeg-linux: check-cgo-enabled
	echo "Installing FFmpeg $(FFMPEG_VERSION) from source..."
	mkdir -p $(FFMPEG_DIRECTORY)
	mkdir -p $(FFMPEG_SRC_DIR)
	cd $(FFMPEG_DIRECTORY) && rm -rf $(FFMPEG_SRC_DIR)
	mkdir -p $(FFMPEG_SRC_DIR)
	echo "Cloning FFmpeg (this may take several minutes)..."
	cd $(FFMPEG_SRC_DIR) && git clone --progress https://github.com/FFmpeg/FFmpeg .
	cd $(FFMPEG_SRC_DIR) && git checkout $(FFMPEG_VERSION)
	cd $(FFMPEG_SRC_DIR) && PKG_CONFIG_PATH="$(X264_DIRECTORY)/lib/pkgconfig:$(OPUS_DIRECTORY)/lib/pkgconfig" ./configure \
		--prefix=$(FFMPEG_DIRECTORY) \
		--enable-gpl \
		--enable-ffplay \
        --enable-libx264 \
        --enable-libopus \
        --enable-alsa \
        --enable-shared \
        --enable-version3 \
        --enable-pic \
        --extra-cflags="-I$(X264_DIRECTORY)/include -I$(OPUS_DIRECTORY)/include" \
		--extra-ldflags="-L$(X264_DIRECTORY)/lib -L$(OPUS_DIRECTORY)/lib"
	cd $(FFMPEG_SRC_DIR) && make -j$(nproc)
	cd $(FFMPEG_SRC_DIR) && make install
	if [ ! -d "$(FFMPEG_DIRECTORY)/lib" ]; then \
		echo "FFmpeg installation failed: lib directory not found"; \
		exit 1; \
	fi

install-ffmpeg-darwin: check-cgo-enabled
	echo "Installing FFmpeg $(FFMPEG_VERSION) from source..."
	mkdir -p $(FFMPEG_DIRECTORY)
	mkdir -p $(FFMPEG_SRC_DIR)
	cd $(FFMPEG_DIRECTORY) && rm -rf $(FFMPEG_SRC_DIR)
	mkdir -p $(FFMPEG_SRC_DIR)
	echo "Cloning FFmpeg (this may take several minutes)..."
	cd $(FFMPEG_SRC_DIR) && git clone --progress https://github.com/FFmpeg/FFmpeg .
	cd $(FFMPEG_SRC_DIR) && git checkout $(FFMPEG_VERSION)
	cd $(FFMPEG_SRC_DIR) && PKG_CONFIG_PATH="$(X264_DIRECTORY)/lib/pkgconfig:$(VPX_DIRECTORY)/lib/pkgconfig" ./configure \
		--prefix=$(FFMPEG_DIRECTORY) \
		--enable-gpl \
		--enable-ffplay \
		--enable-libx264 \
		--enable-libvpx \
		--enable-nonfree \
		--enable-decoder=libvpx_vp8 \
		--enable-decoder=libvpx_vp9 \
		--enable-shared \
		--enable-version3 \
		--enable-pic \
		--extra-cflags="-I$(X264_DIRECTORY)/include -I$(VPX_DIRECTORY)/include" \
		--extra-ldflags="-L$(X264_DIRECTORY)/lib -L$(VPX_DIRECTORY)/lib"
	cd $(FFMPEG_SRC_DIR) && make -j$(nproc)
	cd $(FFMPEG_SRC_DIR) && make install
	if [ ! -d "$(FFMPEG_DIRECTORY)/lib" ]; then \
		echo "FFmpeg installation failed: lib directory not found"; \
		exit 1; \
	fi

install-ffmpeg-docker: check-cgo-enabled
	echo "Installing FFmpeg $(FFMPEG_VERSION) from source..."
	mkdir -p $(FFMPEG_DIRECTORY)
	mkdir -p $(FFMPEG_SRC_DIR)
	cd $(FFMPEG_DIRECTORY) && rm -rf $(FFMPEG_SRC_DIR)
	mkdir -p $(FFMPEG_SRC_DIR)
	echo "Cloning FFmpeg (this may take several minutes)..."
	cd $(FFMPEG_SRC_DIR) && git clone --progress https://github.com/FFmpeg/FFmpeg .
	cd $(FFMPEG_SRC_DIR) && git checkout $(FFMPEG_VERSION)
	cd $(FFMPEG_SRC_DIR) && ./configure \
		--prefix=$(FFMPEG_DIRECTORY) \
		--enable-gpl \
        --enable-shared \
        --enable-version3 \
        --enable-pic
	cd $(FFMPEG_SRC_DIR) && make -j$(nproc)
	cd $(FFMPEG_SRC_DIR) && make install
	if [ ! -d "$(FFMPEG_DIRECTORY)/lib" ]; then \
		echo "FFmpeg installation failed: lib directory not found"; \
		exit 1; \
	fi

install-mavp2p:
	echo "Installing mavp2p from source..."
	mkdir -p $(THIRD_PARTY_DIR)
	git clone https://github.com/bluenviron/mavp2p $(MAVP2P_INSTALL_DIR) 2>$(NULL_DEV) || (cd $(MAVP2P_INSTALL_DIR) && git pull)
	cd $(MAVP2P_INSTALL_DIR) && CGO_ENABLED=0 go build .
	echo "mavp2p installation complete."

# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD
# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

# CGO-dependent build targets
build-fpv-drone: check check-cgo-enabled
	echo "Building fpv drone binary with CGO support..."
	rm -rf $(BUILD_DIR)/fpv/drone
	mkdir -p $(BUILD_DIR)/fpv/drone
	cd $(CMD_DIR)/fpv/drone && \
	$(CGO_RUNTIME_ENV) go build -tags cgo_enabled -o $(BUILD_DIR)/fpv/drone/skyline_sonata.fpv.drone $(LDFLAGS) . || (echo "Build failed"; exit 1)
	echo "fpv drone binary built successfully at $(BUILD_DIR)/fpv/drone"

run-fpv-drone: check-cgo-enabled
	echo "Running fpv drone with CGO support..."
	cd $(BUILD_DIR)/fpv/drone && \
	$(CGO_RUNTIME_ENV) ./skyline_sonata.fpv.drone

# Non-CGO build targets
build-fpv-gcs: check
	echo "Building fpv gcs binary without CGO..."
	rm -rf $(BUILD_DIR)/fpv/gcs
	mkdir -p $(BUILD_DIR)/fpv/gcs
	cd $(CMD_DIR)/fpv/gcs && \
	$(RUNTIME_ENV) CGO_ENABLED=0 go build -o $(BUILD_DIR)/fpv/gcs/skyline_sonata.fpv.gcs $(LDFLAGS) . || (echo "Build failed"; exit 1)
	echo "fpv gcs binary built successfully at $(BUILD_DIR)/fpv/gcs"

run-fpv-gcs:
	echo "Running fpv gcs without CGO..."
	cd $(BUILD_DIR)/fpv/gcs && \
	$(RUNTIME_ENV) ./skyline_sonata.fpv.gcs

build-simple-rtsp-server:
	echo "Building simple rtsp server"
	rm -rf $(BUILD_DIR)/rtsp/server
	mkdir -p $(BUILD_DIR)/rtsp/server
	cd $(WORKING_DIR)/dependencies/services/cmd/rtsp && \
	go build -o $(BUILD_DIR)/rtsp/server/simple-rtsp-server . || (echo "Build failed"; exit 1)
	echo "simple rtsp server binary built successfully at $(BUILD_DIR)/rtsp/server"

run-simple-rtsp-server:
	echo "Running simple-rtsp-server ..."
	cd $(BUILD_DIR)/rtsp/server && \
	./simple-rtsp-server