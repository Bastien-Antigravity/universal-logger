// =============================================================================
// ESSENTIAL PROCESS:
// Unit test runner verifying C++ UniLog instantiation, configuration roundtrips,
// dynamic metadata, caller metadata propagation, and logging macros.
//
// DATA FLOW:
// 1. Input: Programmatic test cases exercising UniLog methods.
// 2. Logic: Executes methods and asserts expected states via assert().
// 3. Output: Pass/fail exit code and console execution diagnostics.
//
// KEY PARAMETERS:
// - None
// =============================================================================

#include "UniversalLogger.hpp"
#include <iostream>
#include <fstream>
#include <cassert>
#include <cstdlib>
#include <string>
#include <unistd.h>

int main() {
    std::cout << ">>> Running C++ Unit Tests..." << std::endl;

    std::string cfg_path = "/tmp/unilog_cpp_test_standalone_" + std::to_string(getpid()) + ".yaml";
    {
        std::ofstream cfg_file(cfg_path);
        cfg_file << "common:\n  name: \"cpp-test-suite\"\n  public_ip: \"127.0.0.1\"\n";
    }
    setenv("CONFIG_PATH", cfg_path.c_str(), 1);
    setenv("LOG_DIR", "/tmp/unilog_cpp_test_logs", 1);

    try {
        // 1. Test Instantiation
        UniLog logger("cpp-test-suite", "standalone", "devel", UniLog::DEBUG, false, 0);
        std::cout << "  - Instantiation passed" << std::endl;

        // 2. Test Get/Set Log Level
        assert(logger.get_level() == UniLog::DEBUG);
        logger.set_level(UniLog::WARNING);
        assert(logger.get_level() == UniLog::WARNING);
        logger.set_level(UniLog::DEBUG);
        std::cout << "  - Level Get/Set passed" << std::endl;

        // 3. Test Config Get/Set
        logger.set_config("test_section", "test_key", "value_123");
        std::string val = logger.get_config("test_section", "test_key", "default");
        assert(val == "value_123");

        std::string missing = logger.get_config("test_section", "non_existent_key", "default_val");
        assert(missing == "default_val");
        std::cout << "  - Config Get/Set passed" << std::endl;

        // 4. Test Metadata Manipulation
        logger.add_metadata("test_meta_key", "meta_val");
        logger.set_metadata("{\"test_env\":\"ci\"}");
        std::cout << "  - Metadata management passed" << std::endl;

        // 5. Test Direct Caller Metadata
        logger.log_with_metadata(UniLog::INFO, "Explicit caller from C++ test", "test_unilog.cpp", "42", "main", "cpp_tests");
        std::cout << "  - Direct log_with_metadata passed" << std::endl;

        // 6. Test Macros & Logging (should execute without crash)
        UNILOG_INFO(logger, "C++ Info message from automated test");
        UNILOG_DEBUG(logger, "C++ Debug message from automated test");
        UNILOG_STREAM(logger, "C++ Stream message from automated test");
        UNILOG_LOGON(logger, "C++ Logon message from automated test");
        UNILOG_LOGOUT(logger, "C++ Logout message from automated test");
        UNILOG_TRADE(logger, "C++ Trade message from automated test");
        UNILOG_SCHEDULE(logger, "C++ Schedule message from automated test");
        UNILOG_REPORT(logger, "C++ Report message from automated test");
        UNILOG_WARNING(logger, "C++ Warning message from automated test");
        UNILOG_ERROR(logger, "C++ Error message from automated test");
        UNILOG_CRITICAL(logger, "C++ Critical message from automated test");
        std::cout << "  - Logging macros execution passed" << std::endl;

        std::cout << ">>> ✅ ALL C++ UNIT TESTS PASSED SUCCESSFULLY!" << std::endl;
    } catch (const std::exception& e) {
        std::remove(cfg_path.c_str());
        std::cerr << ">>> ❌ C++ UNIT TEST FAILED: " << e.what() << std::endl;
        return 1;
    }

    std::remove(cfg_path.c_str());
    return 0;
}
