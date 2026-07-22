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
    } else {
        unilog_dir.join("libunilog")
    };
    
    let target_os = env::var("CARGO_CFG_TARGET_OS").unwrap_or_default();
    let is_windows = target_os == "windows" || (target_os.is_empty() && cfg!(windows));
    let lib_ext = if is_windows { "dll" } else { "so" };
    let lib_name = format!("libunilog.{}", lib_ext);
    let lib_path = lib_dir.join(&lib_name);

    // On Linux/CI, try to build the Go library if missing
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
            }
        } else {
            println!("cargo:warning=Go command not found. Linker may fail.");
        }
    }

    println!("cargo:rustc-link-search=native={}", lib_dir.display());
    println!("cargo:rustc-link-lib=dylib=unilog");
}
