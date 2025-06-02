#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""libgen.py downloads and extracts pre-compiled libcurl-impersonate
libraries for various platforms based on a JSON configuration."""

import json
import shutil
import tarfile
import tempfile
import urllib.request
from pathlib import Path

CURL_IMPERSONATE_VERSION = "1.0.0"
GITHUB_REPO_OWNER = "lexiforest"
GITHUB_REPO_NAME = "curl-impersonate"

PROJECT_ROOT = Path(__file__).parent.parent.resolve()
LIBS_JSON_PATH = PROJECT_ROOT / "misc" / "libs.json"
BASE_OUTPUT_DIR = PROJECT_ROOT / "libs"


def download_file(url, dest_path):
    print(f"Downloading {url} to {dest_path}...")
    try:
        with urllib.request.urlopen(url) as response, open(
            dest_path, "wb"
        ) as out_file:
            shutil.copyfileobj(response, out_file)
        print("Download complete.")
    except Exception as e:
        print(f"Error downloading {url}: {e}")
        raise


def extract_archive(archive_path, extract_to_dir):
    print(f"Extracting {archive_path} to {extract_to_dir}...")
    try:
        with tarfile.open(archive_path, "r:gz") as tar:
            tar.extractall(path=extract_to_dir)
        print("Extraction complete.")
    except Exception as e:
        print(f"Error extracting {archive_path}: {e}")
        raise


def copy_files_from_extraction(extracted_path, dest_path, method_config):
    dest_path.mkdir(parents=True, exist_ok=True)
    method = method_config["extraction_method"]
    print(f"Copying files using method: {method}...")

    if method == "static_linux_bundle":
        copied_any = False
        for pattern in ["lib/*.a", "*.a"]:
            source_dir_glob = extracted_path.glob(pattern)
            for item in source_dir_glob:
                if item.is_file():
                    print(f"  Copying {item.name} to {dest_path}")
                    shutil.copy2(item, dest_path / item.name)
                    copied_any = True
            if copied_any and pattern == "lib/*.a":
                break
        if not copied_any:
            print(
                f"Warning: No .a files found for static_linux_bundle in {extracted_path}"
            )

    elif method == "windows_dll_lib":
        copied_any = False
        bin_dir = extracted_path / "bin"
        target_dll = bin_dir / "libcurl.dll"
        if target_dll.is_file():
            print(f"  Copying {target_dll.name} to {dest_path}")
            shutil.copy2(target_dll, dest_path / target_dll.name)
            copied_any = True
        if not copied_any:
            print(
                f"Warning: 'bin/libcurl.dll' not found for windows_dll_lib in {extracted_path}"
            )

    elif method == "macos_dylib":  # Method name kept, logic changed for .a
        copied_any = False
        # Assuming libcurl-impersonate.a is at the root of extracted_path
        target_static_lib = extracted_path / "libcurl-impersonate.a"
        if target_static_lib.is_file():
            print(f"  Copying {target_static_lib.name} to {dest_path}")
            shutil.copy2(
                target_static_lib, dest_path / target_static_lib.name
            )
            copied_any = True
        if not copied_any:
            print(
                f"Warning: 'libcurl-impersonate.a' not found for macos_dylib (now .a) in {extracted_path}"
            )
    else:
        print(f"Error: Unknown extraction_method '{method}'")
        return
    print("File copying finished.")


def main():
    print(
        "--- Libcurl-Impersonate Library Downloader (All Configured Platforms) ---"
    )
    print(f"Will place libraries into subdirectories of: {BASE_OUTPUT_DIR}")

    if not LIBS_JSON_PATH.exists():
        print(f"Error: Configuration file {LIBS_JSON_PATH} not found.")
        return

    with open(LIBS_JSON_PATH, "r") as f:
        all_configs = json.load(f)

    BASE_OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    for config_index, config in enumerate(all_configs):
        print(
            f"\nProcessing config #{config_index + 1}: {config['match_criteria']}"
        )

        dest_leaf_dir = config["dest_leaf_dir"]
        archive_slug = config["archive_slug"]
        final_lib_dir = BASE_OUTPUT_DIR / dest_leaf_dir

        archive_filename = f"libcurl-impersonate-v{CURL_IMPERSONATE_VERSION}.{archive_slug}.tar.gz"
        download_url = (
            f"https://github.com/{GITHUB_REPO_OWNER}/{GITHUB_REPO_NAME}/releases/download/"
            f"v{CURL_IMPERSONATE_VERSION}/{archive_filename}"
        )

        if final_lib_dir.exists():
            print(f"Cleaning existing libraries directory: {final_lib_dir}")
            try:
                shutil.rmtree(final_lib_dir)
            except OSError as e:
                print(
                    f"Error cleaning directory {final_lib_dir}: {e}. Skipping this platform."
                )
                continue

        final_lib_dir.mkdir(parents=True, exist_ok=True)

        with tempfile.TemporaryDirectory() as tmp_dir_name:
            tmp_dir = Path(tmp_dir_name)
            downloaded_archive_path = tmp_dir / archive_filename
            extracted_contents_path = tmp_dir / "extracted"

            try:
                download_file(download_url, downloaded_archive_path)
                extract_archive(
                    downloaded_archive_path, extracted_contents_path
                )
                copy_files_from_extraction(
                    extracted_contents_path, final_lib_dir, config
                )
                print(
                    f"Successfully downloaded and placed libraries in {final_lib_dir}"
                )
            except Exception as e:
                print(f"An error occurred processing {archive_slug}: {e}")
                print(
                    f"Skipping this platform due to error. Temporary directory {tmp_dir_name} will be cleaned."
                )


if __name__ == "__main__":
    main()