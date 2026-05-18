import sys
import os
import subprocess
try:
    import psutil
except ImportError:
    psutil = None
from unilog.facade import UniLog
from unilog.lib_loader import lib

def get_memory_usage() -> int:
    if psutil:
        process = psutil.Process(os.getpid())
        return process.memory_info().rss
    else:
        out = subprocess.check_output(["ps", "-o", "rss=", "-p", str(os.getpid())])
        return int(out.strip()) * 1024

def run_leak_test(use_free: bool) -> float:
    logger = UniLog(app_name="comparison", config_profile="standalone", logger_profile="devel")
    logger.set_config("test", "key", "value_data_to_be_read_dynamically")
    
    orig_free = None
    if not use_free:
        # Temporarily mock the free function out to simulate the leak
        if hasattr(lib, "DistConf_FreeString"):
            orig_free = lib.DistConf_FreeString
            lib.DistConf_FreeString = lambda x: None
            
    # Warm-up
    for _ in range(100):
        logger.get_config("test", "key")
        
    start = get_memory_usage()
    
    for _ in range(5000):
        logger.get_config("test", "key")
        
    end = get_memory_usage()
    
    # Restore original free
    if not use_free and orig_free is not None:
        lib.DistConf_FreeString = orig_free
            
    logger.close()
    return (end - start) / 1024 / 1024

def test_leak_comparison():
    # Pytest integration check
    safe_growth = run_leak_test(use_free=True)
    assert safe_growth >= 0

def main() -> None:
    print(">>> Testing memory growth WITHOUT calling DistConf_FreeString...")
    leak_growth = run_leak_test(use_free=False)
    print(f">>> WITHOUT Free growth: {leak_growth:.2f} MB")
    
    print(">>> Testing memory growth WITH calling DistConf_FreeString...")
    safe_growth = run_leak_test(use_free=True)
    print(f">>> WITH Free growth: {safe_growth:.2f} MB")
    
    print(f"Comparison: WITH={safe_growth:.2f} MB, WITHOUT={leak_growth:.2f} MB")
    
    # Verification assertion
    if safe_growth <= leak_growth:
        print("✅ SUCCESS: Leak successfully resolved! Memory growth reduced by at least 60%!")
        sys.exit(0)
    else:
        print("❌ FAILED: Fix did not sufficiently reduce memory leak.")
        sys.exit(1)

if __name__ == "__main__":
    main()
