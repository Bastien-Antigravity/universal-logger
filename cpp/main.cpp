#include <iostream>
#include <string>
#include "UniversalLogger.hpp"

int main() {
    try {
        std::cout << ">>> Initializing Universal Logger from C++ Class..." << std::endl;

        // 1. Initialize the logger using RAII
        UniLog logger("standalone", "cpp-app", "standard", UniLog::DEBUG);

        // 2. High-level logging using MACROS (Automatic Metadata!)
        UNILOG_INFO(logger, "Hello from C++ Class with automatic metadata!");
        UNILOG_DEBUG(logger, "Debugging with macros is fast and automatic.");

        // 3. Automated Metadata logging using standardized macros
        UNILOG_WARNING(logger, "System resources running high (detected automatically)!");
        UNILOG_CRITICAL(logger, "Critical failure simulation!");

        // 4. Configuration interaction
        std::string db_ip = logger.get_config("database", "ip", "127.0.0.1 (Default)");
        std::cout << ">>> Config Database IP: " << db_ip << std::endl;

        // 5. Update configuration
        logger.set_config("runtime", "status", "running-cpp");

        std::cout << ">>> Closing session (via Destructor)..." << std::endl;
    } catch (const std::exception& e) {
        std::cerr << "CRITICAL ERROR: " << e.what() << std::endl;
        return 1;
    }

    return 0;
}
