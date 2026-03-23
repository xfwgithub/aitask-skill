#!/usr/bin/env python3
"""
Task Management Skill CLI wrapper

This module provides a Python entry point that delegates to the Go binary.
"""

import json
import os
import subprocess
import sys
from pathlib import Path


def find_go_binary():
    """Find the task-skill Go binary"""
    # Priority 1: Check if running from source (editable install or local development)
    source_binary = Path(__file__).parent.parent / "task-management" / "task-skill"
    if source_binary.exists():
        return str(source_binary)
    
    # Priority 2: Check in bin directory (standard pip install)
    # The Go binary should be installed to the bin directory
    bin_dir = Path(sys.prefix) / "bin"
    binary_path = bin_dir / "task-skill"
    if binary_path.exists() and binary_path.stat().st_size > 100000:  # Should be a large binary (>100KB), not a small script
        return str(binary_path)
    
    # Priority 3: Check in PATH
    for path_dir in os.environ.get("PATH", "").split(os.pathsep):
        binary_path = Path(path_dir) / "task-skill"
        if binary_path.exists() and binary_path.stat().st_size > 100000:
            return str(binary_path)
    
    return None


def main():
    """Main entry point - delegates to Go binary"""
    binary_path = find_go_binary()
    
    if not binary_path:
        print(json.dumps({
            "error": "task-skill binary not found. Please run: pip install task-skill"
        }))
        sys.exit(1)
    
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
