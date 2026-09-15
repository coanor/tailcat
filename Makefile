.DEFAULT_GOAL := build-all

RELEASE_TAGS := $(shell cat build-tags.txt)
VERSION ?= $(shell git describe --tags --always --dirty)
ANDROID_API ?= 24

ifeq ($(shell uname -s),Darwin)
ANDROID_NDK_HOME ?= $(shell brew --prefix 2>/dev/null)/share/android-ndk
ANDROID_NDK_HOST_TAG ?= darwin-x86_64
else
ANDROID_NDK_HOME ?= $(ANDROID_NDK)
ANDROID_NDK_HOST_TAG ?= linux-x86_64
endif

ANDROID_CC ?= $(ANDROID_NDK_HOME)/toolchains/llvm/prebuilt/$(ANDROID_NDK_HOST_TAG)/bin/aarch64-linux-android$(ANDROID_API)-clang

tidy: ## Run go mod tidy and update nix flake hashes
	go mod tidy
	go run ./tool/updateflakes

check-android-cc:
	@test -x "$(ANDROID_CC)" || { \
		echo "Android NDK compiler not found: $(ANDROID_CC)" >&2; \
		echo "Set ANDROID_NDK_HOME or ANDROID_CC, then retry." >&2; \
		exit 1; \
	}

build-android: check-android-cc ## Build the Android ARM64 Termux binary
	mkdir -p dist
	CGO_ENABLED=1 GOOS=android GOARCH=arm64 CC="$(ANDROID_CC)" \
		go build -buildmode=pie -tags "$(RELEASE_TAGS)" \
		-ldflags "-s -w -X main.version=$(VERSION)" \
		-o dist/tailcat-android-arm64 ./cmd/tailcat

build-all: check-android-cc ## Build every GoReleaser binary, including Android ARM64
	ANDROID_CC="$(ANDROID_CC)" goreleaser build --snapshot --clean

release-snapshot: check-android-cc ## Build all release artifacts locally without publishing
	ANDROID_CC="$(ANDROID_CC)" goreleaser release --snapshot --clean

.PHONY: tidy check-android-cc build-android build-all release-snapshot
