# Universal Logger: VBA / Excel Integration

A simple, direct way to integrate high-performance logging and distributed configuration into your Excel or Access applications. Universal Logger uses a Windows-specific wrapper to bridge the Go-shared library into the VBA environment.

## 🚀 Features

- **Standard VBA API**: `Info`, `Debug`, `Warning`, `Error`, and `Critical`.
- **Excel/Access Support**: Works directly from existing `.bas` modules.
- **Background Bridge**: (Planned) Asynchronous updates delivered to the main thread.
- **Easy Deployment**: Requires the `libunilog.dll` to be present on the system.

## 🔧 Installation

1.  **Build the DLL**: (Windows required) `go build -buildmode=c-shared -o libunilog/libunilog.dll src/cgo_bridge/*.go`
2.  **Import the Module**: Import `UniversalLogger.bas` from the `/vba/` directory into your Excel project (Developer Tab -> Visual Basic -> File -> Import File).
3.  **Place the DLL**: Ensure `libunilog.dll` is either in the same folder as your workbook or in a folder in your System PATH.

## 📖 Quick Start

### Basic Logging
```vba
Sub DemoLogging()
    ' Initialize (defaults to standalone config)
    If UniLog_Initialize("standalone", "excel-app") Then
        UniLog_Info "Application has started."
        UniLog_Close
    End If
End Sub
```

### Asynchronous Config Updates (NEW)
To receive real-time configuration updates without crashing Excel, you must start the **Config Watcher**:

```vba
Sub StartMyTool()
    If UniLog_Initialize("production", "my-tool") Then
        ' 1. Start the hidden message pump
        StartConfigWatcher GetUniLogHandle()
        
        ' 2. Updates will now appear in the VBA Immediate Window (Ctrl+G)
        '    or can be handled in UniLog_WindowProc inside the .bas module.
    End If
End Sub

Sub StopMyTool()
    ' 3. Always stop the watcher before closing!
    StopConfigWatcher
    UniLog_Close
End Sub
```

## 🛠️ Configuration and Linking

The `UniversalLogger.bas` file uses `Declare PtrSafe Function` to link with the Go-shared library.

```vba
' Example declaration from .bas
Declare PtrSafe Function UniLog_Init Lib "libunilog.dll" ( ... ) As LongPtr
```

Note: If your DLL is in a custom path, you may need to update the `Lib "libunilog.dll"` line to point to the absolute path of the DLL.

## 🧪 Testing

Refer to [TESTING.md](TESTING.md) for detailed test instructions.
