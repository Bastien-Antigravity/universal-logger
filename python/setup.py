#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import subprocess
import platform
from setuptools import setup, Command
from setuptools.command.build_py import build_py

class BuildGoLib(Command):
    description = "Compile the Go shared library"
    user_options = []

    def initialize_options(self): pass
    def finalize_options(self): pass

    def run(self):
        # 1. Determine platform extension
        ext = ".so"
        if platform.system() == "Darwin":
            ext = ".dylib"
        elif platform.system() == "Windows":
            ext = ".dll"

        target_lib = os.path.join("python", f"libunilog{ext}")
        source_lib = os.path.join("libunilog", f"libunilog{ext}")
        
        # 2. Check if library exists in libunilog/ and copy it to python/ if needed
        if os.path.exists(source_lib):
            print(f">>> Copying {source_lib} to {target_lib} for packaging...")

            import shutil
            shutil.copy2(source_lib, target_lib)
            return

        # 3. Fallback to local Go build if missing
        print(f">>> Missing {source_lib}. Attempting local Go build...")
        cmd = [
            "go", "build", 
            "-buildmode=c-shared",
            "-o", target_lib,
            "./src/cgo_bridge/initialize.go",
            "./src/cgo_bridge/config.go",
            "./src/cgo_bridge/logger.go"
        ]

        
        try:
            # Run from the project root
            subprocess.check_call(cmd, env=os.environ, cwd="..")
        except subprocess.CalledProcessError as e:
            print(f"Failed to build Go library: {e}")
            raise

setup(
    name="unilog",
    version="1.0.0",
    py_modules=["unilog"],
    package_dir={"": "python"},
    package_data={"": ["*.so", "*.dylib", "*.dll", "*.h"]},
    cmdclass={
        "build_go": BuildGoLib,
        "build_py": CustomBuildPy,
    },
    install_requires=[],
    author="Bastien-Antigravity",
    description="Python facade for Universal Logger (Go) - Multi-platform CGO",
)


