# Testing: VBA / Excel Integration

This document explains how to manually verify the VBA facade for Universal Logger.

## Requirements

Before testing in Excel, you must build the Go-shared library as a `.dll` on a Windows machine.

```powershell
# Windows PowerShell build command
go build -buildmode=c-shared -o libunilog/libunilog.dll src/cgo_bridge/*.go
```

## Manual Verification Suite

Testing in VBA is primarily manual. We recommend using the provided `UniversalLogger.bas` in a fresh Excel Workbook.

### 1. Initialization Test
Call `UniLog_Initialize` with the `standalone` profile.
- **Expected Result**: A file log is created (if configured) or the Immediate Window in VBA shows a valid handle value.

### 2. Logging Test
Execute various logging calls:
```vba
UniLog_Info "Test Info"
UniLog_Error "Test Error"
```
- **Expected Result**: Go core should output these logs to the configured sinks (e.g., File, Console, or Network).

### 3. Handle Cleanup Test
Call `UniLog_Close`.
- **Expected Result**: No orphaned `libunilog` processes remain, and the Go session is successfully deleted from the internal `facadeStore`.

## Debugging Tips

### The "Immediate Window"
Use `Debug.Print` in VBA to verify that the library found the DLL and the handle is non-zero.

### Library Pathing Errors
If you receive `Run-time error '53': File not found: libunilog.dll`, ensure the DLL is in the same folder as the `.xlsm` file or in `C:\Windows\System32`.

## Infrastructure Notes

Since VBA is primarily a GUI-driven environment, full automated CI/CD for this component requires a Windows Runner and specialized UI automation. For most releases, manual verification in a standard Excel 64-bit environment is sufficient.
