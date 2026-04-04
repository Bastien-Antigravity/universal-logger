use libc::{c_char, c_int, uintptr_t, free};
use std::ffi::{CStr, CString};
use std::ptr;
use once_cell::sync::Lazy;
use std::sync::Mutex;

// -----------------------------------------------------------------------------
// RAW FFI DECLARATIONS
// -----------------------------------------------------------------------------

#[link(name = "unilog")]
extern "C" {
    fn UniLog_Init(config_profile: *const c_char, app_name: *const c_char, logger_profile: *const c_char, log_level: c_int) -> uintptr_t;
    fn UniLog_Close(handle: uintptr_t);
    fn UniLog_Config_Get(handle: uintptr_t, section: *const c_char, key: *const c_char) -> *mut c_char;
    fn UniLog_Config_Set(handle: uintptr_t, section: *const c_char, key: *const c_char, value: *const c_char);
    fn UniLog_OnMemConfUpdate(handle: uintptr_t, cb: extern "C" fn(*const c_char));
    fn UniLog_LogWithMetadata(handle: uintptr_t, level: i64, msg: *const c_char, file: *const c_char, line: *const c_char, function: *const c_char, module: *const c_char);
    fn UniLog_SetLevel(handle: uintptr_t, level: i64);
}

// -----------------------------------------------------------------------------
// SAFE RUST WRAPPER
// -----------------------------------------------------------------------------

/// Global callback storage for C -> Rust bridge.
/// Currently supports one global callback due to C function pointer limitations.
static GLOBAL_CONFIG_CB: Lazy<Mutex<Option<Box<dyn Fn(String) + Send + 'static>>>> = Lazy::new(|| Mutex::new(None));

extern "C" fn c_callback_bridge(json_data: *const c_char) {
    if let Ok(guard) = GLOBAL_CONFIG_CB.lock() {
        if let Some(ref cb) = *guard {
            let s = unsafe { CStr::from_ptr(json_data).to_string_lossy().into_owned() };
            cb(s);
        }
    }
}

pub enum LogLevel {
    Debug = 1,
    Stream = 2,
    Info = 3,
    Warning = 9,
    Error = 10,
    Critical = 11,
}

pub struct UniLog {
    handle: uintptr_t,
}

impl UniLog {
    /// Initializes a new logger session via the Go shared library.
    pub fn new(config_profile: &str, app_name: &str, logger_profile: &str, log_level: LogLevel) -> Result<Self, String> {
        let c_config = CString::new(config_profile).map_err(|e| e.to_string())?;
        let c_app = CString::new(app_name).map_err(|e| e.to_string())?;
        let c_logger = CString::new(logger_profile).map_err(|e| e.to_string())?;
        
        let handle = unsafe {
            UniLog_Init(c_config.as_ptr(), c_app.as_ptr(), c_logger.as_ptr(), log_level as c_int)
        };
        
        if handle == 0 {
            return Err("Failed to initialize UniLog (Go backend)".to_string());
        }
        
        Ok(UniLog { handle })
    }

    /// Logs a message with custom metadata.
    pub fn log_with_metadata(&self, level: LogLevel, msg: &str, file: &str, line: &str, func: &str, module: &str) {
        let c_msg = CString::new(msg).unwrap_or_default();
        let c_file = CString::new(file).unwrap_or_default();
        let c_line = CString::new(line).unwrap_or_default();
        let c_func = CString::new(func).unwrap_or_default();
        let c_module = CString::new(module).unwrap_or_default();
        
        unsafe {
            UniLog_LogWithMetadata(self.handle, level as i64, c_msg.as_ptr(), c_file.as_ptr(), c_line.as_ptr(), c_func.as_ptr(), c_module.as_ptr());
        }
    }

    /// Convenience logging methods
    pub fn info(&self, msg: &str) { self.log_with_metadata(LogLevel::Info, msg, "lib.rs", "?", "info", "rust-wrapper"); }
    pub fn debug(&self, msg: &str) { self.log_with_metadata(LogLevel::Debug, msg, "lib.rs", "?", "debug", "rust-wrapper"); }
    pub fn warning(&self, msg: &str) { self.log_with_metadata(LogLevel::Warning, msg, "lib.rs", "?", "warning", "rust-wrapper"); }
    pub fn error(&self, msg: &str) { self.log_with_metadata(LogLevel::Error, msg, "lib.rs", "?", "error", "rust-wrapper"); }
    pub fn critical(&self, msg: &str) { self.log_with_metadata(LogLevel::Critical, msg, "lib.rs", "?", "critical", "rust-wrapper"); }

    /// Retrieves a configuration value.
    pub fn get_config(&self, section: &str, key: &str) -> Option<String> {
        let c_sec = CString::new(section).ok()?;
        let c_key = CString::new(key).ok()?;
        
        unsafe {
            let ptr = UniLog_Config_Get(self.handle, c_sec.as_ptr(), c_key.as_ptr());
            if ptr.is_null() {
                return None;
            }
            let res = CStr::from_ptr(ptr).to_string_lossy().into_owned();
            free(ptr as *mut _);
            Some(res)
        }
    }

    /// Updates a configuration value in memory.
    pub fn set_config(&self, section: &str, key: &str, value: &str) {
        let c_sec = CString::new(section).unwrap_or_default();
        let c_key = CString::new(key).unwrap_or_default();
        let c_val = CString::new(value).unwrap_or_default();
        
        unsafe {
            UniLog_Config_Set(self.handle, c_sec.as_ptr(), c_key.as_ptr(), c_val.as_ptr());
        }
    }

    /// Registers a callback to be executed when the configuration is updated in memory.
    pub fn on_config_update<F>(&self, cb: F) 
    where F: Fn(String) + Send + 'static 
    {
        if let Ok(mut guard) = GLOBAL_CONFIG_CB.lock() {
            *guard = Some(Box::new(cb));
        }
        unsafe {
            UniLog_OnMemConfUpdate(self.handle, c_callback_bridge);
        }
    }

    /// Dynamically updates the log level.
    pub fn set_level(&self, level: LogLevel) {
        unsafe {
            UniLog_SetLevel(self.handle, level as i64);
        }
    }
}

impl Drop for UniLog {
    fn drop(&mut self) {
        unsafe {
            UniLog_Close(self.handle);
        }
    }
}

// -----------------------------------------------------------------------------
// LOGGING MACROS
// -----------------------------------------------------------------------------

#[macro_export]
macro_rules! unilog_info {
    ($logger:expr, $msg:expr) => {
        $logger.log_with_metadata($crate::LogLevel::Info, $msg, file!(), &line!().to_string(), "?", module_path!());
    };
}

#[macro_export]
macro_rules! unilog_debug {
    ($logger:expr, $msg:expr) => {
        $logger.log_with_metadata($crate::LogLevel::Debug, $msg, file!(), &line!().to_string(), "?", module_path!());
    };
}

#[macro_export]
macro_rules! unilog_warning {
    ($logger:expr, $msg:expr) => {
        $logger.log_with_metadata($crate::LogLevel::Warning, $msg, file!(), &line!().to_string(), "?", module_path!());
    };
}

#[macro_export]
macro_rules! unilog_critical {
    ($logger:expr, $msg:expr) => {
        $logger.log_with_metadata($crate::LogLevel::Critical, $msg, file!(), &line!().to_string(), "?", module_path!());
    };
}

#[macro_export]
macro_rules! unilog_error {
    ($logger:expr, $msg:expr) => {
        $logger.log_with_metadata($crate::LogLevel::Error, $msg, file!(), &line!().to_string(), "?", module_path!());
    };
}

