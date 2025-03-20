ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    NULL_DEV := NUL
    export PKG_CONFIG := pkgconfiglite
else
    DETECTED_OS := $(shell uname -s)
    NULL_DEV := /dev/null
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

OPUS_DIRECTORY := $(THIRD_PARTY_DIR)/libopus
OPUS_SRC_DIR := $(OPUS_DIRECTORY)/src

# VARIABLES FOR MAVP2P
MAVP2P_INSTALL_DIR := $(THIRD_PARTY_DIR)/mavp2p



# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV-ENV
# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++


COMPILE_ENV := CGO_LDFLAGS="-L$(FFMPEG_DIRECTORY)/lib -L$(X264_DIRECTORY)/lib" \
               CGO_CFLAGS="-I$(FFMPEG_DIRECTORY)/include -I$(X264_DIRECTORY)/include" \
               PKG_CONFIG_PATH="$(FFMPEG_DIRECTORY)/lib/pkgconfig:$(X264_DIRECTORY)/lib/pkgconfig" \
               LD_LIBRARY_PATH="$(FFMPEG_DIRECTORY)/lib:$(X264_DIRECTORY)/lib$LD_LIBRARY_PATH"

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

RUNTIME_ENV := $(HARDWARE_ENV) $(FIREBASE_ENV) $(NETWORK_ENV) $(COMPILE_ENV)


create-env-file:
	@echo "Creating environment files..."
	@# Create runtime environment file
	@echo "# Runtime Environment Variables" > runtime.env
	@echo "$(HARDWARE_ENV)" | tr ' ' '\n' >> runtime.env
	@echo "" >> runtime.env
	@echo "# Firebase Configuration" >> runtime.env
	@echo "$(FIREBASE_ENV)" | tr ' ' '\n' >> runtime.env
	@echo "" >> runtime.env
	@echo "# Network Configuration" >> runtime.env
	@echo "$(NETWORK_ENV)" | tr ' ' '\n' >> runtime.env
	@echo "Runtime environment file created at runtime.env"
	@# Create compilation environment file
	@echo "# Compilation Environment Variables" > compile.env
	@echo "$(COMPILE_ENV)" | tr ' ' '\n' >> compile.env
	@echo "Compilation environment file created at compile.env"
	@# Create complete environment file
	@echo "# Complete Environment Variables" > complete.env
	@echo "$(RUNTIME_ENV)" | tr ' ' '\n' >> complete.env
	@echo "Complete environment file created at complete.env"




# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++


.PHONY: check

check:
	git --version >$(NULL_DEV) 2>&1 || (echo "git is not installed or not in PATH"; exit 1)
	go version >$(NULL_DEV) 2>&1 || (echo "go is not installed or not in PATH"; exit 1)

install-third-party: install-ffmpeg install-mavp2p

install-libx264:
	mkdir -p $(X264_SRC_DIR)
	cd $(X264_SRC_DIR) && git clone https://code.videolan.org/videolan/x264.git .
	cd $(X264_SRC_DIR) && git checkout stable
	cd $(X264_SRC_DIR) && ./configure \
            --prefix=$(X264_DIRECTORY) \
            --enable-shared \
            --enable-pic
	cd $(X264_SRC_DIR) && make -j$(nproc)
	cd $(X264_SRC_DIR) && make install

install-libopus:
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

install-windows-deps:
	if [ "$(DETECTED_OS)" = "Windows" ]; then \
		echo "Installing Windows dependencies..."; \
		pacman -Syu --noconfirm; \
		pacman -S --noconfirm --needed git diffutils mingw-w64-x86_64-toolchain pkg-config make yasm; \
	fi

install-ffmpeg: install-windows-deps
	echo "Installing FFmpeg $(FFMPEG_VERSION) from source..."
	mkdir -p $(FFMPEG_DIRECTORY)
	mkdir -p $(FFMPEG_SRC_DIR)
	cd $(FFMPEG_DIRECTORY) && rm -rf $(FFMPEG_SRC_DIR)
	mkdir -p $(FFMPEG_SRC_DIR)
	echo "Cloning FFmpeg (this may take several minutes)..."
	cd $(FFMPEG_SRC_DIR) && git clone --progress https://github.com/FFmpeg/FFmpeg .
	cd $(FFMPEG_SRC_DIR) && git checkout $(FFMPEG_VERSION)
	cd $(FFMPEG_SRC_DIR) && PKG_CONFIG_PATH="$(X264_DIRECTORY)/lib/pkgconfig" ./configure \
		--prefix=$(FFMPEG_DIRECTORY) \
		--enable-gpl \
		--enable-ffplay \
        --enable-libx264 \
#        --enable-libopus \
#        --enable-alsa \
        --enable-shared \
        --enable-version3 \
        --enable-pic \
        --extra-cflags="-I$(X264_DIRECTORY)/include" \
		--extra-ldflags="-L$(X264_DIRECTORY)/lib"
	cd $(FFMPEG_SRC_DIR) && make -j$(nproc)
	cd $(FFMPEG_SRC_DIR) && make install
	if [ ! -d "$(FFMPEG_DIRECTORY)/lib" ]; then \
		echo "FFmpeg installation failed: lib directory not found"; \
		exit 1; \
	fi
	echo "FFmpeg installation complete. Please set the following environment variables:"
	echo "export CGO_LDFLAGS=\"-L$(FFMPEG_DIRECTORY)/lib/\""
	echo "export CGO_CFLAGS=\"-I$(FFMPEG_DIRECTORY)/include/\""
	echo "export PKG_CONFIG_PATH=\"$(FFMPEG_DIRECTORY)/lib/pkgconfig/\""

install-mavp2p:
	echo "Installing mavp2p from source..."
	mkdir -p $(THIRD_PARTY_DIR)
	git clone https://github.com/bluenviron/mavp2p $(MAVP2P_INSTALL_DIR) 2>$(NULL_DEV) || (cd $(MAVP2P_INSTALL_DIR) && git pull)
	cd $(MAVP2P_INSTALL_DIR) && CGO_ENABLED=0 go build .
	echo "mavp2p installation complete."



# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD-BUILD
# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++


build-delivery-drone: check
	echo "Building delivery drone binary..."
	rm -rf $(BUILD_DIR)/delivery/drone
	mkdir -p $(BUILD_DIR)/delivery/drone
	cd $(CMD_DIR)/delivery/drone && \
	$(RUNTIME_ENV) go build -o $(BUILD_DIR)/delivery/drone/skyline_sonata.delivery.drone $(LDFLAGS) . || (echo "Build failed"; exit 1)
	echo "delivery drone binary built successfully at $(BUILD_DIR)/delivery/drone"

build-delivery-gcs: check
	echo "Building delivery gcs binary..."
	rm -rf $(BUILD_DIR)/delivery/gcs
	mkdir -p $(BUILD_DIR)/delivery/gcs
	cd $(CMD_DIR)/delivery/gcs && \
	$(RUNTIME_ENV) go build -o $(BUILD_DIR)/delivery/gcs/skyline_sonata.delivery.gcs $(LDFLAGS) . || (echo "Build failed"; exit 1)
	echo "delivery gcs binary built successfully at $(BUILD_DIR)/delivery/$(BINARY_GROUND_STATION)"

run-delivery-drone:
	echo "Running delivery drone..."
	cd $(BUILD_DIR)/delivery/drone && \
	$(RUNTIME_ENV) ./skyline_sonata.delivery.drone

run-delivery-gcs:
	echo "Running delivery gcs..."
	cd $(BUILD_DIR)/delivery/gcs && \
	$(RUNTIME_ENV) ./skyline_sonata.delivery.gcs

build-audio-drone: check
	echo "Building audio drone binary..."
	rm -rf $(BUILD_DIR)/audio/drone
	mkdir -p $(BUILD_DIR)/audio/drone
	cd $(CMD_DIR)/audio/drone && \
	$(RUNTIME_ENV) go build -o $(BUILD_DIR)/audio/drone/skyline_sonata.audio.drone $(LDFLAGS) . || (echo "Build failed"; exit 1)
	echo "audio drone binary built successfully at $(BUILD_DIR)/audio/drone"

build-audio-gcs: check
	echo "Building audio gcs binary..."
	rm -rf $(BUILD_DIR)/audio/gcs
	mkdir -p $(BUILD_DIR)/audio/gcs
	cd $(CMD_DIR)/audio/gcs && \
	$(RUNTIME_ENV) go build -o $(BUILD_DIR)/audio/gcs/skyline_sonata.audio.gcs $(LDFLAGS) . || (echo "Build failed"; exit 1)
	echo "audio gcs binary built successfully at $(BUILD_DIR)/audio/$(BINARY_GROUND_STATION)"

run-audio-drone:
	echo "Running audio drone..."
	cd $(BUILD_DIR)/audio/drone && \
	$(RUNTIME_ENV) ./skyline_sonata.audio.drone

run-audio-gcs:
	echo "Running audio gcs..."
	cd $(BUILD_DIR)/audio/gcs && \
	$(RUNTIME_ENV) ./skyline_sonata.audio.gcs