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



COMPILE_ENV := CGO_LDFLAGS="-L$(FFMPEG_DIRECTORY)/lib -L$(X264_DIRECTORY)/lib -L$(VPX_DIRECTORY)/lib" \
               CGO_CFLAGS="-I$(FFMPEG_DIRECTORY)/include -I$(X264_DIRECTORY)/include -I$(VPX_DIRECTORY)/include" \
               PKG_CONFIG_PATH="$(FFMPEG_DIRECTORY)/lib/pkgconfig:$(X264_DIRECTORY)/lib/pkgconfig:$(VPX_DIRECTORY)/lib/pkgconfig" \
               LD_LIBRARY_PATH="$(FFMPEG_DIRECTORY)/lib:$(X264_DIRECTORY)/lib:$(VPX_DIRECTORY)/lib:$$LD_LIBRARY_PATH"

RUNTIME_ENV_MACOS := $(HARDWARE_ENV) $(FIREBASE_ENV) $(NETWORK_ENV) $(COMPILE_ENV) \
                     DYLD_LIBRARY_PATH="$(FFMPEG_DIRECTORY)/lib:$(X264_DIRECTORY)/lib:$(VPX_DIRECTORY)/lib"

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

RUNTIME_ENV := $(HARDWARE_ENV) $(FIREBASE_ENV) $(NETWORK_ENV) $(COMPILE_ENV) $(RUNTIME_ENV_MACOS)

WINDOWS_RUNTIME_ENV := $(HARDWARE_ENV) $(FIREBASE_ENV) $(NETWORK_ENV)


create-env-file:
	@echo "Creating environment files..."
	@# Create runtime environment file
	@echo "# Runtime Environment Variables" > runtime.env
	@echo "$(HARDWARE_ENV)" | sed 's/ \([A-Z_]*=\)/;\1/g' >> runtime.env
	@echo "" >> runtime.env
	@echo "# Firebase Configuration" >> runtime.env
	@echo "$(FIREBASE_ENV)" | sed 's/ \([A-Z_]*=\)/;\1/g' >> runtime.env
	@echo "" >> runtime.env
	@echo "# Network Configuration" >> runtime.env
	@echo "$(NETWORK_ENV)" | sed 's/ \([A-Z_]*=\)/;\1/g' >> runtime.env
	@echo "Runtime environment file created at runtime.env"
	@# Create compilation environment file
	@echo "# Compilation Environment Variables" > compile.env
	@echo "$(COMPILE_ENV)" | sed 's/ \([A-Z_]*=\)/;\1/g' >> compile.env
	@echo "Compilation environment file created at compile.env"




# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET-TARGET
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++


.PHONY: check

check:
	git --version >$(NULL_DEV) 2>&1 || (echo "git is not installed or not in PATH"; exit 1)
	go version >$(NULL_DEV) 2>&1 || (echo "go is not installed or not in PATH"; exit 1)

install-third-party: install-ffmpeg-linux install-mavp2p

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

install-libvpx-darwin:
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

install-ffmpeg-linux:
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

install-ffmpeg-darwin:
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

install-ffmpeg-docker:
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

build-audio-gcs-windows-amd64: check
	echo "Building audio gcs binary..."
	rm -rf $(BUILD_DIR)/delivery/gcs
	mkdir -p $(BUILD_DIR)/delivery/gcs
	cd $(CMD_DIR)/delivery/gcs && \
	$(RUNTIME_ENV) GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/delivery/gcs/skyline_sonata.delivery.windows.gcs $(LDFLAGS) . || (echo "Build failed"; exit 1)
	echo "audio gcs binary built successfully at $(BUILD_DIR)/delivery/$(BINARY_GROUND_STATION)"

run-audio-drone:
	echo "Running audio drone..."
	cd $(BUILD_DIR)/audio/drone && \
	$(RUNTIME_ENV) ./skyline_sonata.audio.drone

run-audio-gcs:
	echo "Running audio gcs..."
	cd $(BUILD_DIR)/audio/gcs && \
	$(RUNTIME_ENV) ./skyline_sonata.audio.gcs

# Network namespace testing targets
run-delivery-drone-with-args:
	echo "Running delivery drone with arguments: $(ARGS)"
	cd $(BUILD_DIR)/delivery/drone && \
	$(RUNTIME_ENV) ./skyline_sonata.delivery.drone $(ARGS)

run-delivery-gcs-with-args:
	echo "Running delivery gcs with arguments: $(ARGS)"
	cd $(BUILD_DIR)/delivery/gcs && \
	$(RUNTIME_ENV) ./skyline_sonata.delivery.gcs $(ARGS)

# Network namespace setup
setup-test-network:
	echo "Setting up test network namespace..."
	# Clean up any existing setup first
	sudo tc qdisc del dev veth0 root 2>/dev/null || true
	sudo ip link del veth0 2>/dev/null || true
	sudo ip netns del testns 2>/dev/null || true
	# Create fresh setup
	sudo ip netns add testns
	sudo ip link add veth0 type veth peer name veth1
	sudo ip link set veth1 netns testns
	sudo ip addr add 192.168.100.1/24 dev veth0 || true
	sudo ip link set veth0 up
	sudo ip netns exec testns ip addr add 192.168.100.2/24 dev veth1 || true
	sudo ip netns exec testns ip link set veth1 up
	sudo ip netns exec testns ip link set lo up
	echo "Test network ready: Main(192.168.100.1) ↔ Namespace(192.168.100.2)"
	echo "Testing connectivity..."
	ping -c 2 192.168.100.2

# Apply network conditions
apply-network-conditions:
	echo "Applying network conditions: $(CONDITIONS)"
	sudo tc qdisc del dev veth0 root 2>/dev/null || true
	sudo tc qdisc add dev veth0 root $(CONDITIONS)
	echo "Applied: $(CONDITIONS)"

# Predefined network conditions
apply-poor-network:
	$(MAKE) apply-network-conditions CONDITIONS="netem delay 200ms loss 3% rate 500kbit"

apply-mobile-network:
	$(MAKE) apply-network-conditions CONDITIONS="netem delay 80ms 20ms loss 1% rate 2mbit"

apply-good-network:
	$(MAKE) apply-network-conditions CONDITIONS="netem delay 50ms loss 0.5% rate 10mbit"

apply-bandwidth-limit:
	$(MAKE) apply-network-conditions CONDITIONS="tbf rate $(RATE) burst 32kbit latency 400ms"

# Run in namespace (server side)
run-delivery-drone-in-namespace:
	echo "Running delivery drone in test namespace (server)..."
	sudo ip netns exec testns bash -c "cd $(BUILD_DIR)/delivery/drone && $(RUNTIME_ENV) ./skyline_sonata.delivery.drone $(ARGS)"

run-delivery-gcs-in-namespace:
	echo "Running delivery gcs in test namespace (server)..."
	sudo ip netns exec testns bash -c "cd $(BUILD_DIR)/delivery/gcs && $(RUNTIME_ENV) ./skyline_sonata.delivery.gcs $(ARGS)"

# Complete test scenarios
test-scenario-bandwidth-degradation:
	echo "=== Testing Bandwidth Degradation Scenario ==="
	$(MAKE) apply-good-network
	echo "Phase 1: Good network (10Mbps) - Run your apps now, press Enter when ready for next phase"
	read
	$(MAKE) apply-mobile-network
	echo "Phase 2: Mobile network (2Mbps) - Press Enter for next phase"
	read
	$(MAKE) apply-poor-network
	echo "Phase 3: Poor network (500kbps) - Press Enter to finish"
	read

test-scenario-packet-loss:
	echo "=== Testing Packet Loss Scenario ==="
	$(MAKE) apply-network-conditions CONDITIONS="netem rate 1mbit loss 1%"
	echo "Phase 1: 1% packet loss - Press Enter for next phase"
	read
	$(MAKE) apply-network-conditions CONDITIONS="netem rate 1mbit loss 5%"
	echo "Phase 2: 5% packet loss - Press Enter to finish"
	read

# Cleanup
cleanup-test-network:
	echo "Cleaning up test network..."
	sudo tc qdisc del dev veth0 root 2>/dev/null || true
	sudo ip link del veth0 2>/dev/null || true
	sudo ip netns del testns 2>/dev/null || true
	echo "Test network cleaned up"

# Helper targets
show-network-status:
	echo "=== Network Status ==="
	sudo tc qdisc show dev veth0 2>/dev/null || echo "No tc rules on veth0"
	sudo ip netns exec testns ip addr show 2>/dev/null || echo "No testns namespace"
	ip addr show veth0 2>/dev/null || echo "No veth0 interface"

# Spawn terminal in namespace
spawn-namespace-terminal:
	echo "Spawning terminal in test namespace..."
	echo "You are now in the network namespace. Run 'exit' to return."
	echo "Your IP in namespace: 192.168.100.2"
	sudo ip netns exec testns bash

# Check required tools
check-tools:
	@echo "Checking required tools for network testing..."
	@command -v ip >/dev/null 2>&1 || { echo "❌ ip command not found (install iproute2)"; exit 1; }
	@command -v tc >/dev/null 2>&1 || { echo "❌ tc command not found (install iproute2)"; exit 1; }
	@command -v ping >/dev/null 2>&1 || { echo "❌ ping command not found"; exit 1; }
	@command -v sudo >/dev/null 2>&1 || { echo "❌ sudo not found"; exit 1; }
	@lsmod | grep -q sch_netem || { echo "⚠️  sch_netem module not loaded (will try to load automatically)"; }
	@echo "✅ All required tools available"

# Help target
help-network-testing:
	@echo "Network Testing Makefile Targets:"
	@echo ""
	@echo "Setup:"
	@echo "  setup-test-network          - Create network namespace and veth pair"
	@echo "  cleanup-test-network        - Remove test network setup"
	@echo ""
	@echo "Run Applications:"
	@echo "  run-delivery-drone-with-args ARGS='--flag=value'"
	@echo "  run-delivery-drone-in-namespace ARGS='--bind=192.168.100.2'"
	@echo ""
	@echo "Network Conditions:"
	@echo "  apply-poor-network          - 200ms delay, 3% loss, 500kbps"
	@echo "  apply-mobile-network        - 80±20ms delay, 1% loss, 2Mbps"
	@echo "  apply-good-network          - 50ms delay, 0.5% loss, 10Mbps"
	@echo "  apply-bandwidth-limit RATE=1mbit"
	@echo ""
	@echo "Complete Test Scenarios:"
	@echo "  test-scenario-bandwidth-degradation"
	@echo "  test-scenario-packet-loss"
	@echo ""
	@echo "Utilities:"
	@echo "  show-network-status         - Show current network state"
	@echo "  test-connection            - Test ping across namespace"
	@echo ""
	@echo "Example Usage:"
	@echo "  make setup-test-network"
	@echo "  make apply-mobile-network"
	@echo "  # Terminal 1:"
	@echo "  make run-delivery-drone-in-namespace ARGS='--bind=192.168.100.2'"
	@echo "  # Terminal 2:"
	@echo "  make run-delivery-gcs-with-args ARGS='--connect=192.168.100.2'"