#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import sys
import json
import time
from unilog import UniLog, LogLevel

def test_metadata_capture():
    print(">>> Starting Python Metadata Capture Test...")
    
    # 1. Initialize UniLog
    # We use 'standalone' and 'devel' for local testing
    logger = UniLog(
        app_name="test-metadata-python",
        config_profile="standalone",
        logger_profile="devel",
        log_level="debug"
    )
    
    try:
        # 2. Add some metadata
        logger.add_metadata("test_id", "python-metadata-123")
        
        # 3. Log messages with automatic metadata capture
        # This will trigger our optimized _get_caller_info using sys._getframe()
        logger.info("Verifying Python metadata capture performance and accuracy")
        logger.debug("Debug log with auto-metadata")
        
        # 4. Manual log with explicit metadata (internal)
        # depth=2 to point to this function
        caller_info = logger._get_caller_info(1)
        print(f"Captured Metadata (sys._getframe):")
        print(f"  File:     {caller_info[0]}")
        print(f"  Line:     {caller_info[1]}")
        print(f"  Function: {caller_info[2]}")
        print(f"  Module:   {caller_info[3]}")
        
        # Verification
        expected_file = os.path.basename(__file__)
        expected_func = "test_metadata_capture"
        
        assert caller_info[0] == expected_file, f"File mismatch: {caller_info[0]} != {expected_file}"
        assert caller_info[2] == expected_func, f"Function mismatch: {caller_info[2]} != {expected_func}"
        assert int(caller_info[1]) > 0, "Line number should be positive"
        
        print("\n✅ Python Metadata Capture verified successfully!")

    finally:
        logger.close()

if __name__ == "__main__":
    # Ensure the unilog package is in the path
    sys.path.append(os.path.join(os.path.dirname(__file__), "unilog"))
    
    try:
        test_metadata_capture()
    except Exception as e:
        print(f"\n❌ Test Failed: {e}")
        sys.exit(1)
