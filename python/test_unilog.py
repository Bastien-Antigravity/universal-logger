#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from unittest import main as unitMain, TestCase as unitTestCase
from unilog import UniLog, LogLevel


##########################################################################
# Logger basic verification

class TestUnilog(unitTestCase):
    
    ##########################################################################
    # Basic logging tests
    
    # Verify that primary logging methods work without crashing using devel profile
    def test_basic_logging(self):
        # We use 'standalone' and 'devel' for testing as they don't require external servers
        with UniLog(config_profile="standalone", app_name="test-basic", logger_profile="devel", log_level="debug") as logger:
            logger.info("Testing Python info log")
            logger.debug("Testing Python debug log")
            logger.warning("Testing Python warning log")
            logger.error("Testing Python error log")
            
            # Test configuration getting (standalone has some defaults)
            val = logger.get_config("logger", "level", "not_found")
            print(f"Logged level from config: {val}")

    # def test_config_updates(self):
    #     # This test is more complex due to the Go-Python callback deadlock potential
    #     # with the Python GIL and synchronous CGO calls.
    #     with UniLog(config_profile="standalone", app_name="test-config") as logger:
    #         updates = []
    #         def on_update(data):
    #             updates.append(data)
    #             print(f"Python received update: {data}")
    #
    #         logger.on_config_update(on_update)
    #         logger.set_config("test_section", "test_key", "test_value")
    #         time.sleep(0.5)
            
    def test_log_level_change(self):
        with UniLog(config_profile="standalone", app_name="test-level", log_level="info") as logger:
            logger.set_level(LogLevel.DEBUG)
            logger.debug("This should be visible after set_level")


##########################################################################
# Entry point

if __name__ == "__main__":
    unitMain()
