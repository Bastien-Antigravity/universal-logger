// =============================================================================
// ESSENTIAL PROCESS:
// Cargo build script for unilog-rs, locating, compiling, and linking libunilog
// shared dynamic library across macOS, Linux, and Windows.
//
// DATA FLOW:
// 1. Input: Target operating system environment variables and project directories.
// 2. Logic: Locates libunilog, triggers Go c-shared build if absent, and emits linker flags.
// 3. Output: Cargo linker directives for dynamic linking against libunilog.
//
// KEY PARAMETERS:
// - LIBUNILOG_PATH / LIBUNILOG_DIR: Optional environment override paths.
// =============================================================================

use std::env;
use std::path::PathBuf;
use std::process::Command;

fn main() {
    let project_dir = env::var("CARGO_MANIFEST_DIR").unwrap();
    let mut unilog_dir = PathBuf::from(&project_dir);
    unilog_dir.pop(); // Go to unilog dir
    
    let mut root_dir = unilog_dir.clone();
    root_dir.pop(); // Go to true universal-logger root

    let lib_dir = if let Ok(val) = env::var("LIBUNILOG_PATH") {
        let p = PathBuf::from(val);
        if p.is_file() {
            p.parent().unwrap().to_path_buf()
        } else {
            p
        }
    } else if let Ok(val) = env::var("LIBUNILOG_DIR") {
        PathBuf::from(val)
    } else if root_dir.join("libunilog").exists() {
        root_dir.join("libunilog")
    } else {
        unilog_dir.join("libunilog")
    };
    
    let target_os = env::var("CARGO_CFG_TARGET_OS").unwrap_or_default();
    let is_windows = target_os == "windows" || (target_os.is_empty() && cfg!(windows));
    let is_macos = target_os == "macos" || (target_os.is_empty() && cfg!(target_os = "macos"));
    let lib_ext = if is_windows {
        "dll"
    } else if is_macos {
        "dylib"
    } else {
        "so"
    };
    let lib_name = format!("libunilog.{}", lib_ext);
    let lib_path = lib_dir.join(&lib_name);

    // On Linux/macOS/CI, try to build the Go library if missing
    let is_ci = env::var("GITHUB_ACTIONS").is_ok();
    if !is_windows && (is_ci || !lib_path.exists()) {
        println!("cargo:warning=Attempting to build Go shared library (cgo_bridge)...");
        let go_src = root_dir.join("src").join("cgo_bridge");
        
        let status = Command::new("go")
            .args(&[
                "build",
                "-buildmode=c-shared",
                "-o",
                lib_path.to_str().unwrap(),
            ])
            .arg(go_src)
            .status();

        if let Ok(s) = status {
            if !s.success() {
                println!("cargo:warning=Go build failed. Linker may fail.");
            } else if is_macos {
                let _ = Command::new("install_name_tool")
                    .args(&["-id", "@rpath/libunilog.dylib", lib_path.to_str().unwrap()])
                    .status();
            }
        } else {
            println!("cargo:warning=Go command not found. Linker may fail.");
        }
    }

    println!("cargo:rustc-link-search=native={}", lib_dir.display());
    println!("cargo:rustc-link-lib=dylib=unilog");
    if is_macos || !is_windows {
        println!("cargo:rustc-link-arg=-Wl,-rpath,{}", lib_dir.display());
    }
}
