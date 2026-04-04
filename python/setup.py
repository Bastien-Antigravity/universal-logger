#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from os import path as osPath, environ as osEnviron
from subprocess import check_call as subprocessCheckCall, CalledProcessError as subprocessCalledProcessError
from platform import system as platformSystem
from shutil import copy2 as shutilCopy2
from setuptools import setup as setuptoolsSetup, Command as setuptoolsCommand


##########################################################################
# Custom Build Commands

class BuildGoLib(setuptoolsCommand):
    """
    Compile the Go shared library for the current platform.
    Used during installation or source distribution builds.
    """
    description = "Compile the Go shared library"
    user_options = []

    def initialize_options(self): pass
    def finalize_options(self): pass

    def run(self):
        # 1. Determine platform extension
        ext = ".so"
        if platformSystem() == "Darwin":
            ext = ".dylib"
        elif platformSystem() == "Windows":
            ext = ".dll"

        target_lib = osPath.join("python", f"libunilog{ext}")
        source_lib = osPath.join("libunilog", f"libunilog{ext}")
        
        # 2. Check if library exists in libunilog/ and copy it to python/ if needed
        if osPath.exists(source_lib):
            print(f">>> Copying {source_lib} to {target_lib} for packaging...")
            shutilCopy2(source_lib, target_lib)
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
            # Run from the project root (.. relative to python dir)
            subprocessCheckCall(cmd, env=osEnviron, cwd="..")
        except subprocessCalledProcessError as e:
            print(f"Failed to build Go library: {e}")
            raise


##########################################################################
# Package configuration

setuptoolsSetup(
    name="unilog",
    version="1.0.0",
    packages=["unilog"],
    package_dir={"": "python"},
    package_data={"unilog": ["*.so", "*.dylib", "*.dll", "*.h"]},
    # Custom build commands including the Go bridge compilation
    cmdclass={
        "build_go": BuildGoLib,
    },
    install_requires=[],
    author="Bastien-Antigravity",
    description="Python facade for Universal Logger (Go) - Multi-platform CGO",
)
