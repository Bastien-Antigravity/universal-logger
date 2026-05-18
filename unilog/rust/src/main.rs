use unilog_rs::{UniLog, LogLevel, unilog_info, unilog_debug, unilog_warning};

fn main() {
    println!(">>> Initializing Universal Logger from Rust...");

    // 1. Initialize the logger safe wrapper (Order: app_name, config_profile, logger_profile)
    let logger = match UniLog::new("rust-app-demo", "standalone", "devel", LogLevel::Debug, false, 0) {
        Ok(l) => l,
        Err(e) => {
            eprintln!("Error initializing logger: {}", e);
            return;
        }
    };

    // 2. Metadata management
    logger.add_metadata("version", "1.1.0");
    logger.add_metadata("language", "rust");
    
    let mut meta = std::collections::HashMap::new();
    meta.insert("env".to_string(), "development".to_string());
    meta.insert("arch".to_string(), "parity-hardened".to_string());
    logger.set_metadata(&meta);

    // 3. Register a configuration update callback
    logger.on_config_update(|json_data| {
        println!(">>> [CALLBACK] Config Updated: {}", json_data);
    });

    // 4. Verify Log Level
    println!(">>> Initial Log Level: {:?}", logger.get_level() as i32);

    // 5. High-level logging using MACROS (Automatic Metadata!)
    unilog_info!(logger, "Hello from safe Rust bindings with auto-metadata!");
    unilog_debug!(logger, "Debugging with macros is zero-cost and automatic.");

    // 6. Automated Metadata logging (Warning)
    unilog_warning!(logger, "System resources running high (detected automatically)!");

    // 7. Domain-specific levels
    logger.log_with_metadata(LogLevel::Trade, "Simulated trade execution", file!(), &line!().to_string(), "main", module_path!());

    // 8. Configuration interaction (Explicit names)
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
