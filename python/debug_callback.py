#!/usr/bin/env python3
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
