import sys
import os
try:
    import psutil
except ImportError:
    psutil = None
from unilog.facade import UniLog

def get_memory_usage() -> int:
    if psutil:
        process = psutil.Process(os.getpid())
        return process.memory_info().rss
    else:
        # Fallback for systems without psutil (uses standard macOS shell command ps)
        import subprocess
        out = subprocess.check_output(["ps", "-o", "rss=", "-p", str(os.getpid())])
        return int(out.strip()) * 1024  # ps RSS is in KB, convert to bytes

def main() -> None:
    print(">>> Starting Python FFI memory leak verification test...")
    # Initialize the logger
    logger = UniLog(app_name="memory-leak-test", config_profile="standalone", logger_profile="devel")
    
    # Baseline configuration key-value
    logger.set_config("test_mem", "leak_key", "leak_value_data")
    
    # 1. Warm-up
    for _ in range(100):
        logger.get_config("test_mem", "leak_key")
        
    start_mem = get_memory_usage()
    print(f">>> Initial memory after warm-up: {start_mem / 1024 / 1024:.2f} MB")
    
    # 2. Loop calling get_config
    iterations = 50000
    for i in range(iterations):
        val = logger.get_config("test_mem", "leak_key")
        if val != "leak_value_data":
            print(f"Error: unexpected value retrieved: {val}")
            sys.exit(1)
            
        if (i + 1) % 10000 == 0:
            current_mem = get_memory_usage()
            print(f">>> Iteration {i+1}/{iterations} memory: {current_mem / 1024 / 1024:.2f} MB")
            
    end_mem = get_memory_usage()
    growth = end_mem - start_mem
    growth_mb = growth / 1024 / 1024
    print(f">>> Final memory: {end_mem / 1024 / 1024:.2f} MB")
    print(f">>> Memory growth: {growth_mb:.2f} MB")
    
    logger.close()
    
    # If there is a leak, 50,000 C allocations will leak several MB.
    # Without a leak, growth should be extremely close to 0 MB.
    if growth_mb > 1.2:
        print("❌ FAILED: Significant memory leak detected!")
        sys.exit(1)
    else:
        print("✅ SUCCESS: Zero-leak FFI verified successfully!")
        sys.exit(0)

if __name__ == "__main__":
    main()
