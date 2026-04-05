#ifndef UNIVERSAL_LOGGER_HPP
#define UNIVERSAL_LOGGER_HPP

#include <string>
#include <stdexcept>
#include <iostream>
#include "../libunilog/libunilog.h"

/**
 * @class UniLog
 * @brief C++ RAII Wrapper for the Universal Logger (Go) shared library.
 * 
 * Provides an idiomatic C++ interface with automated resource management.
 */
class UniLog {
public:
    enum Level {
        DEBUG = 1,
        STREAM = 2,
        INFO = 3,
        WARNING = 9,
        ERROR = 10,
        CRITICAL = 11
    };

    /**
     * @brief Initialize a new logger session.
     */
    UniLog(const std::string& config_profile = "standalone", 
                    const std::string& app_name = "cpp-app", 
                    const std::string& logger_profile = "standard", 
                    int log_level = INFO,
                    bool use_local_notifier = false) {
        
        handle_ = UniLog_Init(
            const_cast<char*>(config_profile.c_str()),
            const_cast<char*>(app_name.c_str()),
            const_cast<char*>(logger_profile.c_str()),
            log_level,
            use_local_notifier ? 1 : 0
        );

        if (handle_ == 0) {
            throw std::runtime_error("Failed to initialize Universal Logger Go backend.");
        }
    }

    /**
     * @brief Cleanup the logger session.
     */
    ~UniLog() {
        if (handle_ != 0) {
            UniLog_Close(handle_);
            handle_ = 0;
        }
    }

    // Disable copying to prevent double-close of the handle
    UniLog(const UniLog&) = delete;
    UniLog& operator=(const UniLog&) = delete;

    // --- Logging Methods ---

    void log(int level, const std::string& msg, 
             const std::string& file = "unknown", 
             const std::string& line = "0", 
             const std::string& func = "unknown", 
             const std::string& module = "cpp") {
        
        UniLog_LogWithMetadata(
            handle_, level, 
            const_cast<char*>(msg.c_str()), 
            const_cast<char*>(file.c_str()), 
            const_cast<char*>(line.c_str()), 
            const_cast<char*>(func.c_str()), 
            const_cast<char*>(module.c_str())
        );
    }

    void debug(const std::string& msg) { log(DEBUG, msg); }
    void info(const std::string& msg) { log(INFO, msg); }
    void warning(const std::string& msg) { log(WARNING, msg); }
    void error(const std::string& msg) { log(ERROR, msg); }
    void critical(const std::string& msg) { log(CRITICAL, msg); }

    // --- Configuration Methods ---

    /**
     * @brief Get a configuration value.
     */
    std::string get_config(const std::string& section, const std::string& key, const std::string& default_val = "") {
        char* val = UniLog_Config_Get(handle_, const_cast<char*>(section.c_str()), const_cast<char*>(key.c_str()));
        if (!val) {
            return default_val;
        }
        std::string result(val);
        free(val); // CGO returns a C.CString which must be freed
        return result;
    }

    /**
     * @brief Update a configuration value in memory.
     */
    void set_config(const std::string& section, const std::string& key, const std::string& value) {
        UniLog_Config_Set(
            handle_, 
            const_cast<char*>(section.c_str()), 
            const_cast<char*>(key.c_str()), 
            const_cast<char*>(value.c_str())
        );
    }

    /**
     * @brief Dynamically update log level.
     */
    void set_level(int level) {
        UniLog_SetLevel(handle_, level);
    }

    /**
     * @brief Set a callback for notifications.
     */
    void set_notification_callback(void (*cb)(const char*)) {
        UniLog_RegisterNotifCallback(handle_, cb);
    }

private:
    GoUintptr handle_;
};

// -----------------------------------------------------------------------------
// LOGGING MACROS (for automatic metadata capture)
// -----------------------------------------------------------------------------

#define UNILOG_INFO(logger, msg) \
    (logger).log(UniLog::INFO, (msg), __FILE__, std::to_string(__LINE__), __FUNCTION__, "cpp-module")

#define UNILOG_DEBUG(logger, msg) \
    (logger).log(UniLog::DEBUG, (msg), __FILE__, std::to_string(__LINE__), __FUNCTION__, "cpp-module")

#define UNILOG_WARNING(logger, msg) \
    (logger).log(UniLog::WARNING, (msg), __FILE__, std::to_string(__LINE__), __FUNCTION__, "cpp-module")

#define UNILOG_ERROR(logger, msg) \
    (logger).log(UniLog::ERROR, (msg), __FILE__, std::to_string(__LINE__), __FUNCTION__, "cpp-module")

#define UNILOG_CRITICAL(logger, msg) \
    (logger).log(UniLog::CRITICAL, (msg), __FILE__, std::to_string(__LINE__), __FUNCTION__, "cpp-module")


#endif // UNIVERSAL_LOGGER_HPP

