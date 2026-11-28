#!/usr/bin/env python
# coding:utf-8
"""
=============================================================================
ESSENTIAL PROCESS:
Package setup and distribution script for the unilog Python client library.

DATA FLOW:
1. Input: Repo VERSION.txt, Python source modules, and native libunilog shared library.
2. Logic: Validates version fail-fast and bundles unilog modules with CGO binaries/headers.
3. Output: Self-contained Python distribution wheel/package.

KEY PARAMETERS:
- name: Package distribution name ("unilog").
=============================================================================
"""

from os import environ as osEnviron
from os.path import exists as osPathExists, join as osPathJoin
from pathlib import Path as pathlibPath
from platform import system as platformSystem
from setuptools import Command as setuptoolsCommand, setup as setuptoolsSetup
from shutil import copy2 as shutilCopy2
from subprocess import CalledProcessError as subprocessCalledProcessError, check_call as subprocessCheckCall


# -----------------------------------------------------------------------------------------------
# Custom Build Commands

class BuildGoLib(setuptoolsCommand):
    """
    Compile or copy the Go shared library for the current platform into unilog/.
    """

    Name = "BuildGoLib"
    description = "Compile the Go shared library"
    user_options = []

    # -----------------------------------------------------------------------------------------------

    def initialize_options(self) -> None:
        pass

    # -----------------------------------------------------------------------------------------------

    def finalize_options(self) -> None:
        pass

    # -----------------------------------------------------------------------------------------------

    def run(self) -> None:
        ext = ".so"
        if platformSystem() == "Darwin":
            ext = ".dylib"
        elif platformSystem() == "Windows":
            ext = ".dll"

        target_lib = osPathJoin("unilog", f"libunilog{ext}")
        target_hdr = osPathJoin("unilog", "libunilog.h")
        source_lib = osPathJoin("..", "..", "libunilog", f"libunilog{ext}")
        source_hdr = osPathJoin("..", "..", "libunilog", "libunilog.h")

        if osPathExists(source_lib):
            print(f">>> Copying {source_lib} to {target_lib} for packaging...")
            shutilCopy2(source_lib, target_lib)
            if osPathExists(source_hdr):
                shutilCopy2(source_hdr, target_hdr)
            return

        print(f">>> Missing {source_lib}. Attempting local Go build...")
        cmd = [
            "go",
            "build",
            "-buildmode=c-shared",
            "-o",
            osPathJoin("unilog", "python", "unilog", f"libunilog{ext}"),
            "./src/cgo_bridge",
        ]

        try:
            subprocessCheckCall(cmd, env=osEnviron, cwd=osPathJoin("..", ".."))
        except subprocessCalledProcessError as err:
            print(f"Failed to build Go library: {err}")
            raise


# -----------------------------------------------------------------------------------------------

def _load_version() -> str:
    repo_root = pathlibPath(__file__).resolve().parent.parent.parent
    version_path = repo_root / "VERSION.txt"
    if not version_path.is_file():
        raise FileNotFoundError(f"Missing VERSION.txt at {version_path}")
    version = version_path.read_text(encoding="utf-8").strip()
    if not version:
        raise ValueError(f"Empty VERSION.txt at {version_path}")
    return version


# -----------------------------------------------------------------------------------------------
# Package Configuration

setuptoolsSetup(
    name="unilog",
    version=_load_version(),
    packages=["unilog"],
    package_dir={"": "."},
    package_data={"unilog": ["*.so", "*.dylib", "*.dll", "*.h"]},
    include_package_data=True,
    cmdclass={
        "build_go": BuildGoLib,
    },
    install_requires=[],
    author="Bastien-Antigravity",
    description="Python facade for Universal Logger (Go) - Multi-platform CGO",
)
