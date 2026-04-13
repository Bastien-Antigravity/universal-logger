use std::env;
use std::path::PathBuf;

fn main() {
    let project_dir = env::var("CARGO_MANIFEST_DIR").unwrap();
    let mut lib_dir = PathBuf::from(project_dir);
    // Path to libunilog folder relative to universal-logger/rust/
    lib_dir.pop(); 
    lib_dir.push("libunilog");

    println!("cargo:rustc-link-search=native={}", lib_dir.display());
    
    // We generated libunilog.a (MinGW style)
    // On Windows with MSVC, we might need to rename it or use a .lib
    // But let's try to link it as 'unilog'
    println!("cargo:rustc-link-lib=dylib=unilog");
}
