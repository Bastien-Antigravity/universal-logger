# Universal Logger: VBA / Excel Integration

A simple, direct way to integrate high-performance logging and distributed configuration into your Excel or Access applications. Universal Logger uses a Windows-specific wrapper to bridge the Go-shared library into the VBA environment.

## 🚀 Native Core (DLL)

Since Modernization 2026, the VBA integration relies on the **compiled Go shared library** (`libunilog.dll`).

**Requirement**: Ensure `libunilog.dll` is in the same directory as your workbook (`.xlsm`, `.xlsb`) or in your system `PATH`.

## 🚀 Features

- **Standard VBA API**: `GetConfig`, `SetConfig`, and `UniLog_LogWithMetadata`.
- **Excel/Access Support**: Works directly from existing `.bas` modules with 64-bit and 32-bit compatibility.
- **Background Bridge**: Asynchronous updates delivered via a hidden Windows Message Pump for zero-crash stability.
- **Telemetry Parity**: Shared library ensures your Excel logs match the format and performance of Go/Rust/Python services.

## 🔧 Installation

1.  **Build the DLL**: (Windows required) See root README for Docker-based build instructions.
2.  **Import the Module**: Import `UniversalLogger.bas` from the `/vba/` directory into your Excel project (Developer Tab -> Visual Basic -> File -> Import File).
3.  **Place the DLL**: Place `libunilog.dll` next to your workbook.

## 📖 Quick Start

### Basic Logging
```vba
Sub DemoLogging()
    Dim handle As LongPtr
    
    ' 1. Initialize (loads libunilog.dll)
    handle = UniLog_Init("standalone", "Excel-Tool", "standard", Level_INFO, 0)
    
    If handle <> 0 Then
        ' 2. Log message
        UniLog_LogWithMetadata handle, Level_INFO, "Application has started from VBA", "Module1.bas", "10", "DemoLogging", "Excel-VBA"
        
        ' 3. Clean up
        UniLog_Close handle
    End If
End Sub
```

### Asynchronous Config Updates
To receive real-time configuration updates without locking up the Excel UI:

```vba
Sub StartMyTool()
    Dim h As LongPtr
    h = UniLog_Init("production", "my-tool", "standard", Level_INFO, 1)
    
    If h <> 0 Then
        ' 1. Start the hidden message pump
        StartConfigWatcher h
        
        ' 2. Updates will appear in the Immediate Window (Ctrl+G)
        '    Dispatcher logic in UniLog_WindowProc inside the .bas module.
    End If
End Sub

Sub StopMyTool()
    ' 3. Always stop the watcher before closing!
    StopConfigWatcher
    ' UniLog_Close (handle) - store your handle globally if needed
End Sub
```

## 🛠️ Configuration and Linking

The `UniversalLogger.bas` file uses `Declare PtrSafe Function` to link with the Go-shared library.

```vba
' Example declaration from .bas
Declare PtrSafe Function UniLog_Init Lib "libunilog.dll" ( ... ) As LongPtr
```

## 🧪 Testing

Refer to [TESTING.md](TESTING.md) for detailed test instructions.
