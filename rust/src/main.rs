use unilog_rs::{UniLog, LogLevel, unilog_info, unilog_debug, unilog_warning};

fn main() {
    println!(">>> Initializing Universal Logger from Rust...");

    // 1. Initialize the logger safe wrapper
    let logger = match UniLog::new("standalone", "rust-app-debug", "standard", LogLevel::Debug) {
        Ok(l) => l,
        Err(e) => {
            eprintln!("Error initializing logger: {}", e);
            return;
        }
    };

    // 2. Register a configuration update callback
    logger.on_config_update(|json_data| {
        println!(">>> [CALLBACK] Config Updated: {}", json_data);
    });

    // 3. High-level logging using MACROS (Automatic Metadata!)
    unilog_info!(logger, "Hello from safe Rust bindings with auto-metadata!");
    unilog_debug!(logger, "Debugging with macros is zero-cost and automatic.");

    // 4. Automated Metadata logging (Warning)
    unilog_warning!(logger, "System resources running high (detected automatically)!");

    // 5. Configuration interaction (Explicit names)
    if let Some(db_ip) = logger.get_config("database", "ip") {
        println!(">>> Config Database IP: {}", db_ip);
    }

    // 6. Update configuration (triggers callback)
    println!(">>> Updating runtime status (should trigger callback)...");
    logger.set_config("runtime", "status", "running-rust-auto");

    // Give some time for the callback to run (it runs in a Go-managed thread)
    std::thread::sleep(std::time::Duration::from_millis(500));

    println!(">>> Closing session (via Drop trait)...");
}
