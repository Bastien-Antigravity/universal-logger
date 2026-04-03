import os
import subprocess
import platform
import shutil
from setuptools import setup, Command
from setuptools.command.build_py import build_py

class BuildGoLib(Command):
    description = "Compile the Go shared library"
    user_options = []

    def initialize_options(self): pass
    def finalize_options(self): pass

    def run(self):
        # 1. Resolve library name based on platform
        ext = ".so"
        if platform.system() == "Darwin":
            ext = ".dylib"
        elif platform.system() == "Windows":
            ext = ".dll"

        target_lib = os.path.join("python", f"libdistconf_flexlog{ext}")
        
        # 2. Run go build
        print(f"Building Go shared library: {target_lib}")
        cmd = [
            "go", "build", 
            "-buildmode=c-shared",
            "-o", target_lib,
            "src/facade/facade.go",
            "src/facade/c_api.go"
        ]
        
        try:
            subprocess.check_call(cmd, env=os.environ, cwd="..")
        except subprocess.CalledProcessError as e:
            print(f"Failed to build Go library: {e}")
            raise

class CustomBuildPy(build_py):
    def run(self):
        self.run_command("build_go")
        super().run()

setup(
    name="distconf-flexlog",
    version="1.0.0",
    packages=["distconf_flexlog"],
    package_dir={"": "python"},
    package_data={"distconf_flexlog": ["*.so", "*.dylib", "*.dll"]},
    cmdclass={
        "build_go": BuildGoLib,
        "build_py": CustomBuildPy,
    },
    install_requires=[],
    author="Bastien-Antigravity",
    description="Python facade for Distributed Config and Flexible Logger",
)
