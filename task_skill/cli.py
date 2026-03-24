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


VERSION = "1.1.6"
GITHUB_RELEASES_URL = f"https://github.com/xfwgithub/aitask-skill/releases/download/v{VERSION}"


def find_go_binary():
    """Find the task-skill Go binary"""
    # Priority 1: Check if running from source (editable install or local development)
    source_binary = Path(__file__).parent.parent / "task-management" / "task-skill"
    if source_binary.exists():
        if check_binary_version(source_binary):
            return str(source_binary)
    
    # Priority 2: Check in package data directory
    package_dir = Path(__file__).parent
    binary_path = package_dir / "task-skill"
    if binary_path.exists() and binary_path.stat().st_size > 100000:
        if check_binary_version(binary_path):
            return str(binary_path)
    
    # Priority 3: Check in bin directory (standard pip install)
    bin_dir = Path(sys.prefix) / "bin"
    binary_path = bin_dir / "task-skill"
    if binary_path.exists() and binary_path.stat().st_size > 100000:
        if check_binary_version(binary_path):
            return str(binary_path)
    
    # Priority 4: Check in PATH
    for path_dir in os.environ.get("PATH", "").split(os.pathsep):
        binary_path = Path(path_dir) / "task-skill"
        if binary_path.exists() and binary_path.stat().st_size > 100000:
            if check_binary_version(binary_path):
                return str(binary_path)
    
    return None

def check_binary_version(binary_path):
    """Check if the binary version matches the expected version"""
    try:
        result = subprocess.run([str(binary_path), "--version"], capture_output=True, text=True, timeout=2)
        if result.returncode == 0:
            # We relax the strict version check to avoid infinite download loops 
            # when the Github Release package is slightly out of sync with the pip package.
            output = result.stdout.strip()
            # Still warn if version is completely wrong or very old, but accept it if it runs
            if VERSION not in output:
                print(f"Warning: Found binary ({output}), pip wrapper expects {VERSION}.", file=sys.stderr)
            return True
        return False
    except Exception:
        return False


def get_package_dir():
    """Get the task-skill package directory (where static resources should be)"""
    # For pip installs, use the package directory where cli.py is located
    package_dir = Path(__file__).parent
    return package_dir


def download_binary():
    """Download the task-skill package from GitHub Releases"""
    system = platform.system()
    machine = platform.machine()
    
    if system != "Darwin":
        print(json.dumps({
            "error": f"Unsupported platform: {system} {machine}. Only macOS is supported."
        }))
        sys.exit(1)
    
    # Determine architecture
    if machine == "arm64":
        arch = "arm64"
    else:
        print(json.dumps({
            "error": f"Unsupported architecture: {machine}. Only ARM64 (Apple Silicon) is supported."
        }))
        sys.exit(1)
    
    # Download the full zip package (includes binary + static resources)
    zip_name = f"task-skill-v{VERSION}.zip"
    download_url = f"{GITHUB_RELEASES_URL}/{zip_name}"
    
    # Determine where to extract
    package_dir = get_package_dir()
    if not package_dir:
        print(json.dumps({
            "error": "Could not determine package directory"
        }))
        sys.exit(1)
    
    print(f"Downloading task-skill package from GitHub Releases...", file=sys.stderr)
    print(f"URL: {download_url}", file=sys.stderr)
    print(f"Destination: {package_dir}", file=sys.stderr)
    
    try:
        # Download the zip file
        import tempfile
        import zipfile
        import shutil
        
        with tempfile.NamedTemporaryFile(suffix='.zip', delete=False) as tmp_file:
            tmp_zip_path = Path(tmp_file.name)
            with urllib.request.urlopen(download_url) as response:
                tmp_zip_path.write_bytes(response.read())
        
        print(f"✓ Package downloaded ({tmp_zip_path.stat().st_size:,} bytes)", file=sys.stderr)
        
        # Extract to temp directory first
        with tempfile.TemporaryDirectory() as tmp_extract_dir:
            with zipfile.ZipFile(str(tmp_zip_path), 'r') as zip_ref:
                zip_ref.extractall(tmp_extract_dir)
            
            # The zip contains a parent directory (task-skill-vX.X.X/)
            # Find and copy its contents to package_dir
            extracted_dirs = list(Path(tmp_extract_dir).iterdir())
            if not extracted_dirs:
                raise Exception("No directories found in zip")
            
            source_dir = extracted_dirs[0]  # task-skill-vX.X.X/
            
            # Copy all files to package directory
            for item in source_dir.iterdir():
                dest = package_dir / item.name
                if dest.exists():
                    if dest.is_dir():
                        shutil.rmtree(dest)
                    else:
                        dest.unlink()
                
                if item.is_dir():
                    shutil.copytree(item, dest)
                else:
                    shutil.copy2(item, dest)
        
        # Clean up temp file
        tmp_zip_path.unlink()
        
        print(f"✓ Package extracted successfully", file=sys.stderr)
        
        # Return path to binary
        binary_path = package_dir / "task-skill"
        if binary_path.exists():
            binary_path.chmod(0o755)
            return str(binary_path)
        else:
            print(json.dumps({
                "error": "Binary not found after extraction"
            }))
            sys.exit(1)
            
    except Exception as e:
        print(json.dumps({
            "error": f"Failed to download/extract package: {str(e)}",
            "help": f"Please download manually from: {download_url}"
        }))
        sys.exit(1)


def main():
    """Main entry point - delegates to Go binary"""
    binary_path = find_go_binary()
    
    # If binary not found, try to download it from GitHub Releases
    if not binary_path:
        binary_path = download_binary()
    
    # Check if we need to run in server mode
    server_mode = "--server" in sys.argv[1:]
    
    # For server mode, change working directory to package directory
    # so Go binary can find templates/ and static/ directories
    cwd = None
    if server_mode:
        cwd = str(Path(binary_path).parent)
    
    # Set unified database path to prevent split brain between CLI and server mode
    env = os.environ.copy()
    if "TASK_SKILL_DB_PATH" not in env:
        db_dir = Path.home() / ".task-skill"
        db_dir.mkdir(parents=True, exist_ok=True)
        env["TASK_SKILL_DB_PATH"] = str(db_dir / "tasks.db")
    
    # Pass all arguments to the Go binary
    try:
        result = subprocess.run(
            [binary_path] + sys.argv[1:],
            capture_output=True,
            text=True,
            cwd=cwd,
            env=env
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
