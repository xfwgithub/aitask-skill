#!/usr/bin/env python3
"""
Task Management Skill - Python wrapper for Go binary
Makes the Go CLI installable via pip
"""

import os
import platform
import subprocess
import sys
from pathlib import Path

from setuptools import setup, find_packages
from setuptools.command.install import install
from setuptools.command.develop import develop

# Package metadata
NAME = "task-skill"
VERSION = "0.4.0"
DESCRIPTION = "零依赖、高性能的任务管理技能 - Task Management Skill"
AUTHOR = "xfwgithub"
URL = "https://github.com/xfwgithub/aitask-skill"

# Go binary configuration
GO_BINARY_NAME = "task-skill"
GO_SOURCE_DIR = "task-management"


class BuildGoBinary:
    """Build the Go binary during installation"""
    
    @staticmethod
    def check_go_installed():
        """Check if Go is installed"""
        try:
            result = subprocess.run(
                ["go", "version"],
                capture_output=True,
                text=True,
                check=True
            )
            return True, result.stdout.strip()
        except (subprocess.CalledProcessError, FileNotFoundError):
            return False, None
    
    @staticmethod
    def build_go_binary(source_dir, output_dir):
        """Build the Go binary"""
        # Check platform
        system = platform.system()
        machine = platform.machine()
        
        if system != "Darwin" or machine not in ["arm64", "x86_64"]:
            print(f"Warning: This package is designed for macOS Apple Silicon (arm64).")
            print(f"Current platform: {system} {machine}")
            print("The binary may not work correctly on this platform.")
        
        # Check Go
        has_go, go_version = BuildGoBinary.check_go_installed()
        if not has_go:
            raise RuntimeError(
                "Go is not installed. Please install Go 1.21+ from https://golang.org/dl/"
            )
        
        print(f"Found: {go_version}")
        
        # Build binary
        source_path = Path(source_dir).resolve()
        if not source_path.exists():
            raise RuntimeError(f"Source directory not found: {source_path}")
        
        output_path = Path(output_dir) / GO_BINARY_NAME
        
        print(f"Building Go binary...")
        print(f"Source: {source_path}")
        print(f"Output: {output_path}")
        
        env = os.environ.copy()
        env["GOOS"] = "darwin"
        env["GOARCH"] = "arm64"
        env["CGO_ENABLED"] = "1"
        
        result = subprocess.run(
            ["go", "build", "-o", str(output_path), "."],
            cwd=str(source_path),
            capture_output=True,
            text=True,
            env=env
        )
        
        if result.returncode != 0:
            raise RuntimeError(f"Failed to build Go binary:\n{result.stderr}")
        
        # Make executable
        output_path.chmod(0o755)
        
        print(f"✓ Built: {output_path}")
        return str(output_path)


class CustomInstallCommand(install):
    """Custom install command that builds Go binary"""
    
    def run(self):
        # Run standard install
        install.run(self)
        
        # Build Go binary
        print("\n" + "="*60)
        print("Building Task Management Skill (Go binary)...")
        print("="*60)
        
        # Find source directory
        source_dir = Path(__file__).parent / GO_SOURCE_DIR
        
        # Build to the installation scripts directory
        scripts_dir = Path(self.install_scripts)
        scripts_dir.mkdir(parents=True, exist_ok=True)
        
        try:
            binary_path = BuildGoBinary.build_go_binary(source_dir, scripts_dir)
            print(f"\n✓ Task-skill installed to: {binary_path}")
            print(f"✓ You can now use 'task-skill' command\n")
        except RuntimeError as e:
            print(f"\n✗ Build failed: {e}")
            print("\nYou can manually build the binary:")
            print(f"  cd {source_dir}")
            print("  bash build.sh")
            raise


class CustomDevelopCommand(develop):
    """Custom develop command for editable install"""
    
    def run(self):
        develop.run(self)
        
        print("\n" + "="*60)
        print("Building Task Management Skill (Go binary)...")
        print("="*60)
        
        source_dir = Path(__file__).parent / GO_SOURCE_DIR
        
        # For editable install, build to the source scripts directory
        try:
            binary_path = BuildGoBinary.build_go_binary(source_dir, source_dir)
            print(f"\n✓ Task-skill built: {binary_path}")
            print(f"✓ Add to PATH: export PATH=$PATH:{source_dir}\n")
        except RuntimeError as e:
            print(f"\n✗ Build failed: {e}")
            raise


# Read README for long description
readme_path = Path(__file__).parent / "README.md"
long_description = ""
if readme_path.exists():
    long_description = readme_path.read_text(encoding="utf-8")


setup(
    name=NAME,
    version=VERSION,
    description=DESCRIPTION,
    long_description=long_description,
    long_description_content_type="text/markdown",
    author=AUTHOR,
    url=URL,
    packages=find_packages(),
    package_data={
        "": ["*.md", "*.go", "*.mod", "*.sum", "build.sh", "start.sh"],
    },
    include_package_data=True,
    python_requires=">=3.10",
    cmdclass={
        "install": CustomInstallCommand,
        "develop": CustomDevelopCommand,
    },
    entry_points={
        "console_scripts": [
            "task-skill=task_skill.cli:main",
        ],
    },
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Operating System :: MacOS :: MacOS X",
        "Programming Language :: Go",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        "Topic :: Software Development :: Libraries :: Python Modules",
        "Topic :: Utilities",
    ],
    keywords="task management, cli, agent, ai, skill",
    project_urls={
        "Bug Reports": f"{URL}/issues",
        "Source": URL,
    },
)
