import asyncio
import unittest
from unilog import UniLog, LogLevel

class TestUnifiedCallback(unittest.TestCase):
    def test_sync_callback(self):
        print("\n>>> Testing Synchronous Callback...")
        updates = []
        def my_cb(data):
            print(f"Sync Callback received: {data}")
            updates.append(data)

        with UniLog(config_profile="standalone", app_name="test-sync", logger_profile="devel") as logger:
            logger.on_config_update(my_cb)
            
            # Trigger update
            logger.set_config("test", "key", "value1")
            
            # Give a small amount of time for the background Go thread to call back
            import time
            time.sleep(0.2)
            
        self.assertGreater(len(updates), 0)
        self.assertEqual(updates[-1]["test"]["key"], "value1")

    def test_async_iterator(self):
        print("\n>>> Testing Async Iterator...")
        
        async def run_test():
            async with UniLog(config_profile="standalone", app_name="test-async", logger_profile="devel") as logger:
                # Get the listener
                listener = logger.on_config_update()
                
                # We'll run the listener in a task so we can trigger updates
                received = []
                
                async def listen_task():
                    async for update in listener:
                        print(f"Async Iterator received: {update}")
                        received.append(update)
                        if len(received) >= 1:
                            break
                
                t = asyncio.create_task(listen_task())
                
                # Trigger update
                await asyncio.sleep(0.1)
                logger.set_config("async_test", "status", "active")
                
                # Wait for listener to finish
                await asyncio.wait_for(t, timeout=2.0)
                
                self.assertEqual(len(received), 1)
                self.assertEqual(received[0]["async_test"]["status"], "active")

        asyncio.run(run_test())

    def test_dual_mode(self):
        print("\n>>> Testing Dual Mode (Sync + Async simultaneously)...")
        
        async def run_test():
            sync_received = []
            def sync_cb(data):
                sync_received.append(data)

            async with UniLog(config_profile="standalone", app_name="test-dual", logger_profile="devel") as logger:
                # Register Sync
                logger.on_config_update(sync_cb)
                
                # Start Async
                received = []
                async def listen_task():
                    async for update in logger.on_config_update():
                        received.append(update)
                        if len(received) >= 1:
                            break
                
                t = asyncio.create_task(listen_task())
                
                # Trigger update
                await asyncio.sleep(0.1)
                logger.set_config("dual", "mode", "enabled")
                
                await asyncio.wait_for(t, timeout=2.0)
                
                # Wait a bit for sync (though it should be immediate)
                await asyncio.sleep(0.1)
                
                self.assertEqual(len(received), 1)
                self.assertEqual(len(sync_received), 1)

        asyncio.run(run_test())

if __name__ == "__main__":
    unittest.main()
