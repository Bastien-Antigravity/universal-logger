import time
from facade import DistconfFlexlogFacade

def main():
    print("Initializing Distconf-Flexlog Facade from Python...")
    
    # 1. Initialize the facade
    # Use 'devel' profile for easy console output
    facade = DistconfFlexlogFacade(
        config_profile="standalone", 
        app_name="python-tester", 
        logger_profile="devel", 
        log_level="debug"
    )

    # 2. Register a callback for config updates
    def on_update(data):
        print(f"\n[Python Callback] Config Updated: {data}")

    facade.on_config_update(on_update)

    # 3. Log some messages
    print("Logging messages with automatic stack inspection...")
    facade.debug("This is a DEBUG message from Python")
    facade.info("This is an INFO message from Python")
    facade.warning("This is a WARNING message from Python")
    facade.error("This is an ERROR message from Python")

    # 4. Wait a bit to ensure async operations (if any) are done
    time.sleep(1)
    
    print("\nCleanup...")
    facade.close()

if __name__ == "__main__":
    main()
