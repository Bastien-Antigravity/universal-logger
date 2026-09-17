#!/usr/bin/env python3
"""
=============================================================================
ESSENTIAL PROCESS:
Diagnostic debug utility for interactively observing FFI callback behavior
and runtime thread transitions.

DATA FLOW:
1. Input: Interactive configuration mutation events.
2. Logic: Prints detailed diagnostics when CGO callbacks fire.
3. Output: Console trace of callback execution and payload decoding.

KEY PARAMETERS:
- None
=============================================================================
"""

import sys
import os
import time

# Ensure we pick up the new package structure
sys.path.insert(0, os.path.abspath("python"))

from unilog import UniLog, LogLevel

def test_callback():
    print("--- Starting Callback Test ---")
    
    received_updates = []
    def my_cb(data):
        print(f"!!! Python: CB RECEIVED: {data}")
        received_updates.append(data)

    print("!!! Python: Initializing UniLog...")
    logger = UniLog()
    
    print("!!! Python: Registering Callback...")
    logger.on_config_update(my_cb)
    
    print("!!! Python: Triggering Update via set_config...")
    logger.set_config("test_section", "test_key", "test_value")
    
    print("!!! Python: Waiting for callback...")
    for _ in range(10):
        if received_updates:
            break
        time.sleep(0.5)
        print("...still waiting...")

    if received_updates:
        print("SUCCESS: Callback received updates!")
    else:
        print("FAILURE: No updates received.")
    
    logger.close()

if __name__ == "__main__":
    test_callback()
