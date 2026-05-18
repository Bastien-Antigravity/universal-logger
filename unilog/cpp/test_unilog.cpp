#include "UniversalLogger.hpp"
#include <iostream>
#include <cassert>
#include <string>

int main() {
    std::cout << ">>> Running C++ Unit Tests..." << std::endl;

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

        // 5. Test Macros & Logging (should execute without crash)
        UNILOG_INFO(logger, "C++ Info message from automated test");
        UNILOG_DEBUG(logger, "C++ Debug message from automated test");
        UNILOG_WARNING(logger, "C++ Warning message from automated test");
        UNILOG_ERROR(logger, "C++ Error message from automated test");
        UNILOG_CRITICAL(logger, "C++ Critical message from automated test");
        std::cout << "  - Logging macros execution passed" << std::endl;

        std::cout << ">>> ✅ ALL C++ UNIT TESTS PASSED SUCCESSFULLY!" << std::endl;
    } catch (const std::exception& e) {
        std::cerr << ">>> ❌ C++ UNIT TEST FAILED: " << e.what() << std::endl;
        return 1;
    }

    return 0;
}
