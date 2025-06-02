.ONESHELL:
SHELL := bash

VERSION	:= 1.0.0
CURL_VERSION := curl-8_13_0

BUILD_DIR := build
GO_CURL_PACKAGE_DIR := .
HEADERS_DEST_DIR := $(GO_CURL_PACKAGE_DIR)/libs/include/curl

CURL_ORIGINAL_SRC_DIR := $(BUILD_DIR)/curl_original_source_for_patching
CURL_IMPERSONATE_SRC_DIR := $(BUILD_DIR)/curl_impersonate_source

.PHONY: preprocess_headers
preprocess_headers: $(BUILD_DIR)/.headers_preprocessed

$(BUILD_DIR)/.headers_preprocessed: $(CURL_ORIGINAL_SRC_DIR)/.patched_and_ready_for_headers
	@echo "==> Copying patched headers to $(HEADERS_DEST_DIR)..."
	mkdir -p $(HEADERS_DEST_DIR)
	cp -R $(CURL_ORIGINAL_SRC_DIR)/include/curl/* $(HEADERS_DEST_DIR)/
	touch $(BUILD_DIR)/.headers_preprocessed
	@echo "==> Header preprocessing complete."

$(CURL_ORIGINAL_SRC_DIR)/.patched_and_ready_for_headers: $(CURL_ORIGINAL_SRC_DIR)/.unpacked_and_ready $(CURL_IMPERSONATE_SRC_DIR)/.unpacked_and_ready
	@echo "==> Patching $(CURL_ORIGINAL_SRC_DIR) with $(CURL_IMPERSONATE_SRC_DIR) patches..."
	cd $(CURL_ORIGINAL_SRC_DIR)
	patch -p1 < ../../$(CURL_IMPERSONATE_SRC_DIR)/patches/curl.patch
	cd ../..
	touch $@

$(CURL_ORIGINAL_SRC_DIR)/.unpacked_and_ready: $(BUILD_DIR)
	@echo "==> Downloading and extracting original libcurl $(CURL_VERSION)..."
	curl -L "https://github.com/curl/curl/archive/refs/tags/$(CURL_VERSION).zip" -o "$(BUILD_DIR)/curl-$(CURL_VERSION)-download.zip"
	rm -rf "$(BUILD_DIR)/curl-$(CURL_VERSION)" $(CURL_ORIGINAL_SRC_DIR)
	unzip -q -o "$(BUILD_DIR)/curl-$(CURL_VERSION)-download.zip" -d $(BUILD_DIR)
	mv "$(BUILD_DIR)/curl-$(CURL_VERSION)" $(CURL_ORIGINAL_SRC_DIR)
	touch $@

$(CURL_IMPERSONATE_SRC_DIR)/.unpacked_and_ready: $(BUILD_DIR)
	@echo "==> Downloading and extracting curl-impersonate $(VERSION) (for patches)..."
	curl -L "https://github.com/lexiforest/curl-impersonate/archive/refs/tags/v$(VERSION).tar.gz" \
		-o "$(BUILD_DIR)/curl-impersonate-$(VERSION)-download.tar.gz"
	rm -rf "$(BUILD_DIR)/curl-impersonate-$(VERSION)" $(CURL_IMPERSONATE_SRC_DIR)
	tar -xzf "$(BUILD_DIR)/curl-impersonate-$(VERSION)-download.tar.gz" -C $(BUILD_DIR)
	mv "$(BUILD_DIR)/curl-impersonate-$(VERSION)" $(CURL_IMPERSONATE_SRC_DIR)
	if [ ! -f "$(CURL_IMPERSONATE_SRC_DIR)/patches/curl.patch" ]; then \
		echo "Error: Patch file $(CURL_IMPERSONATE_SRC_DIR)/patches/curl.patch not found after moving to $(CURL_IMPERSONATE_SRC_DIR)."; \
		echo "Please check actual extracted directory name and VERSION variable."; \
		exit 1; \
	fi
	touch $@

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PHONY: setup_headers
setup_headers: preprocess_headers
	@echo "==> Patched libcurl headers are ready in $(HEADERS_DEST_DIR)."

.PHONY: clean_headers
clean_headers:
	@echo "==> Cleaning up header generation files..."
	rm -rf $(BUILD_DIR)
	rm -rf $(HEADERS_DEST_DIR) 
