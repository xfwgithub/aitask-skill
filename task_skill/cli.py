#!/usr/bin/env python3
"""
Task Management Skill CLI wrapper

This module provides a Python entry point that delegates to the Go binary.
Auto-downloads the binary from GitHub Releases if not found.
"""

import json
import os
import platform
import subprocess
import sys
import urllib.request
from pathlib import Path


VERSION = "0.4.0"
GITHUB_RELEASES_URL = f"https://github.com/xfwgithub/aitask-skill/releases/download/v{VERSION}"


def find_go_binary():
    """Find the task-skill Go binary"""
    # Priority 1: Check if running from source (editable install or local development)
    source_binary = Path(__file__).parent.parent / "task-management" / "task-skill"
    if source_binary.exists():
        return str(source_binary)
    
    # Priority 2: Check in package data directory
    package_dir = Path(__file__).parent
    binary_path = package_dir / "task-skill"
    if binary_path.exists() and binary_path.stat().st_size > 100000:
        return str(binary_path)
    
    # Priority 3: Check in bin directory (standard pip install)
    bin_dir = Path(sys.prefix) / "bin"
    binary_path = bin_dir / "task-skill"
    if binary_path.exists() and binary_path.stat().st_size > 100000:
        return str(binary_path)
    
    # Priority 4: Check in PATH
    for path_dir in os.environ.get("PATH", "").split(os.pathsep):
        binary_path = Path(path_dir) / "task-skill"
        if binary_path.exists() and binary_path.stat().st_size > 100000:
            return str(binary_path)
    
    return None


def download_binary():
    """Download the Go binary from GitHub Releases"""
    system = platform.system()
    machine = platform.machine()
    
    if system != "Darwin":
        print(json.dumps({
            "error": f"Unsupported platform: {system} {machine}. Only macOS is supported."
        }))
        sys.exit(1)
    
    # Determine binary name
    if machine == "arm64":
        binary_name = "task-skill"  # ARM64 binary
    else:
        print(json.dumps({
            "error": f"Unsupported architecture: {machine}. Only ARM64 (Apple Silicon) is supported."
        }))
        sys.exit(1)
    
    # Download URL
    download_url = f"{GITHUB_RELEASES_URL}/{binary_name}"
    
    # Determine where to save the binary
    package_dir = Path(__file__).parent
    binary_path = package_dir / "task-skill"
    
    print(f"Downloading task-skill binary from GitHub Releases...", file=sys.stderr)
    print(f"URL: {download_url}", file=sys.stderr)
    print(f"Destination: {binary_path}", file=sys.stderr)
    
    try:
        # Download the binary
        with urllib.request.urlopen(download_url) as response:
            binary_data = response.read()
        
        # Save to package directory
        binary_path.write_bytes(binary_data)
        binary_path.chmod(0o755)
        
        print(f"✓ Binary downloaded successfully ({len(binary_data):,} bytes)", file=sys.stderr)
        return str(binary_path)
        
    except Exception as e:
        print(json.dumps({
            "error": f"Failed to download binary: {str(e)}",
            "help": f"Please download manually from: {GITHUB_RELEASES_URL}"
        }))
        sys.exit(1)


def main():
    """Main entry point - delegates to Go binary"""
    binary_path = find_go_binary()
    
    # If binary not found, try to download it from GitHub Releases
    if not binary_path:
        binary_path = download_binary()
    
    # Pass all arguments to the Go binary
    try:
        result = subprocess.run(
            [binary_path] + sys.argv[1:],
            capture_output=True,
            text=True
        )
        
        # Print stdout (JSON output from Go binary)
        if result.stdout:
            print(result.stdout, end="")
        
        # Print stderr if any
        if result.stderr:
            print(result.stderr, file=sys.stderr)
        
        sys.exit(result.returncode)
        
    except Exception as e:
        print(json.dumps({
            "error": f"Failed to run task-skill: {str(e)}"
        }))
        sys.exit(1)


if __name__ == "__main__":
    main()
