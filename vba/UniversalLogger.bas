Attribute VB_Name = "UniLog"
' -----------------------------------------------------------------------------
' UniLog - VBA Bridge for Microsoft Excel
' -----------------------------------------------------------------------------
' This module allows calling the Universal Logger (Go) shared library from VBA.
' Ensure that libunilog.dll is in the same directory as the workbook or in a
' directory listed in the system PATH.
' -----------------------------------------------------------------------------

Option Explicit

    ' 64-bit Excel
    Private Declare PtrSafe Function UniLog_Init Lib "libunilog.dll" (ByVal configProfile As String, ByVal appName As String, ByVal loggerProfile As String, ByVal logLevel As Long) As LongPtr
    Private Declare PtrSafe Sub UniLog_Close Lib "libunilog.dll" (ByVal handle As LongPtr)
    Private Declare PtrSafe Function UniLog_Config_Get Lib "libunilog.dll" (ByVal handle As LongPtr, ByVal section As String, ByVal key As String) As LongPtr
    Private Declare PtrSafe Sub UniLog_Config_Set Lib "libunilog.dll" (ByVal handle As LongPtr, ByVal section As String, ByVal key As String, ByVal value As String)
    Private Declare PtrSafe Sub UniLog_LogWithMetadata Lib "libunilog.dll" (ByVal handle As LongPtr, ByVal level As LongLong, ByVal msg As String, ByVal file As String, ByVal line As String, ByVal functionName As String, ByVal moduleName As String)
    Private Declare PtrSafe Sub UniLog_SetLevel Lib "libunilog.dll" (ByVal handle As LongPtr, ByVal level As LongLong)
    
    ' --- VBA CALLBACK BRIDGE (NEW) ---
    Private Declare PtrSafe Sub UniLog_RegisterVBAWindow Lib "libunilog.dll" (ByVal handle As LongPtr, ByVal hwnd As LongPtr, ByVal msgId As Long)

    ' Windows API for Message Pump
    Private Declare PtrSafe Function CreateWindowExA Lib "user32" (ByVal dwExStyle As Long, ByVal lpClassName As String, ByVal lpWindowName As String, ByVal dwStyle As Long, ByVal x As Long, ByVal y As Long, ByVal nWidth As Long, ByVal nHeight As Long, ByVal hWndParent As LongPtr, ByVal hMenu As LongPtr, ByVal hInstance As LongPtr, lpParam As Any) As LongPtr
    Private Declare PtrSafe Function DestroyWindow Lib "user32" (ByVal hwnd As LongPtr) As Long
    Private Declare PtrSafe Function SetWindowLongPtrA Lib "user32" Alias "SetWindowLongPtrA" (ByVal hwnd As LongPtr, ByVal nIndex As Long, ByVal dwNewLong As LongPtr) As LongPtr
    Private Declare PtrSafe Function CallWindowProcA Lib "user32" (ByVal lpPrevWndFunc As LongPtr, ByVal hwnd As LongPtr, ByVal msg As Long, ByVal wParam As LongPtr, ByVal lParam As LongPtr) As LongPtr
    Private Declare PtrSafe Function GetModuleHandleA Lib "kernel32" (ByVal lpModuleName As String) As LongPtr

    ' Constants
    Private Const GWLP_WNDPROC As Long = -4
    Private Const WM_USER As Long = &H400
    Private Const HWND_MESSAGE As Long = -3
    
    ' User-definable Message ID for the callback
    Private Const UNILOG_UPDATE_MSG As Long = WM_USER + 101

    ' Private state for the message pump
    Private hProxyWnd As LongPtr
    Private pOldWndProc As LongPtr
    Private hUniLogHandle As LongPtr

    ' Helper to convert C string pointer to VBA String
    Private Declare PtrSafe Function lstrlenA Lib "kernel32" (ByVal lpString As LongPtr) As Long
    Private Declare PtrSafe Sub lstrcpyA Lib "kernel32" (ByVal lpString1 As String, ByVal lpString2 As LongPtr)
#Else
    ' 32-bit Excel (Note: libunilog.dll must also be 32-bit)
    Private Declare Function UniLog_Init Lib "libunilog.dll" (ByVal configProfile As String, ByVal appName As String, ByVal loggerProfile As String, ByVal logLevel As Long) As Long
    Private Declare Sub UniLog_Close Lib "libunilog.dll" (ByVal handle As Long)
    Private Declare Function UniLog_Config_Get Lib "libunilog.dll" (ByVal handle As Long, ByVal section As String, ByVal key As String) As Long
    Private Declare Sub UniLog_Config_Set Lib "libunilog.dll" (ByVal handle As Long, ByVal section As String, ByVal key As String, ByVal value As String)
    Private Declare Sub UniLog_LogWithMetadata Lib "libunilog.dll" (ByVal handle As Long, ByVal level As Long, ByVal msg As String, ByVal file As String, ByVal line As String, ByVal functionName As String, ByVal moduleName As String)
    Private Declare Sub UniLog_SetLevel Lib "libunilog.dll" (ByVal handle As Long, ByVal level As Long)
#End If

' Shared Log Levels
Public Enum UniLogLevel
    Level_DEBUG = 1
    Level_STREAM = 2
    Level_INFO = 3
    Level_WARNING = 9
    Level_ERROR = 10
    Level_CRITICAL = 11
End Enum

' Helper to convert C-string pointer to VBA string
Private Function PtrToString(ByVal ptr As LongPtr) As String
    Dim length As Long
    Dim res As String
    If ptr = 0 Then
        PtrToString = ""
        Exit Function
    End If
    length = lstrlenA(ptr)
    res = Space$(length)
    lstrcpyA res, ptr
    PtrToString = res
End Function

' -----------------------------------------------------------------------------
' VBA MESSAGE PUMP IMPLEMENTATION
' -----------------------------------------------------------------------------

' The Window Procedure that receives messages from the Go core.
' MUST be public and in a standard module to use AddressOf.
Public Function UniLog_WindowProc(ByVal hwnd As LongPtr, ByVal msg As Long, ByVal wParam As LongPtr, ByVal lParam As LongPtr) As LongPtr
    On Error Resume Next ' Preventive check 
    
    If msg = UNILOG_UPDATE_MSG Then
        ' lParam contains the pointer to the JSON string from Go
        Dim jsonUpdate As String
        jsonUpdate = PtrToString(lParam)
        
        ' --- DISPATCH TO USER CALLBACK ---
        ' For simplicity in this facade, we print to Immediate Window
        ' but users can replace this with a call to their own function.
        Debug.Print "!!! Universal Logger Update: " & jsonUpdate
        
        ' Note: In a real-world scenario, you might want to call a 
        ' specific 'public sub' defined by the user here.
        
        UniLog_WindowProc = 0
        Exit Function
    End If
    
    ' Pass all other messages to the original window procedure
    UniLog_WindowProc = CallWindowProcA(pOldWndProc, hwnd, msg, wParam, lParam)
End Function

' Initializes the hidden window and registers it with Go
Public Sub StartConfigWatcher(ByVal handle As LongPtr)
    If hProxyWnd <> 0 Then Exit Sub ' Already running
    hUniLogHandle = handle
    
    ' 1. Create a "Message-Only" window (invisible, no taskbar)
    hProxyWnd = CreateWindowExA(0, "Static", "UniLogProxy", 0, 0, 0, 0, 0, HWND_MESSAGE, 0, GetModuleHandleA(vbNullString), ByVal 0&)
    
    If hProxyWnd = 0 Then
        Debug.Print "!!! Error: Could not create UniLog Proxy Window"
        Exit Sub
    End If
    
    ' 2. Subclass the window to use our UniLog_WindowProc
    pOldWndProc = SetWindowLongPtrA(hProxyWnd, GWLP_WNDPROC, AddressOf UniLog_WindowProc)
    
    ' 3. Tell the Go DLL about our Window and the Message ID we expect
    UniLog_RegisterVBAWindow handle, hProxyWnd, UNILOG_UPDATE_MSG
    
    Debug.Print "!!! VBA: Config Watcher Started (HWND: " & hProxyWnd & ")"
End Sub

' Cleans up the hidden window
Public Sub StopConfigWatcher()
    If hProxyWnd = 0 Then Exit Sub
    
    ' 1. Restore the original window procedure
    SetWindowLongPtrA hProxyWnd, GWLP_WNDPROC, pOldWndProc
    
    ' 2. Destroy the window
    DestroyWindow hProxyWnd
    
    hProxyWnd = 0
    pOldWndProc = 0
    Debug.Print "!!! VBA: Config Watcher Stopped."
End Sub

' -----------------------------------------------------------------------------
' VBA MESSAGE PUMP IMPLEMENTATION (Windows Only)
' -----------------------------------------------------------------------------

' The Window Procedure that receives messages from the Go core.
' MUST be public and in a standard module to use AddressOf.
Public Function UniLog_WindowProc(ByVal hwnd As LongPtr, ByVal msg As Long, ByVal wParam As LongPtr, ByVal lParam As LongPtr) As LongPtr
    On Error Resume Next 
    
    If msg = UNILOG_UPDATE_MSG Then
        ' lParam contains the pointer to the JSON string from Go
        Dim jsonUpdate As String
        jsonUpdate = PtrToString(lParam)
        
        ' Dispatch to the active listener if any
        ' (Users can add their own event dispatching logic here)
        Debug.Print "!!! UniLog Update Received: " & jsonUpdate
        
        UniLog_WindowProc = 0
        Exit Function
    End If
    
    ' Pass all other messages to the original window procedure
    UniLog_WindowProc = CallWindowProcA(pOldWndProc, hwnd, msg, wParam, lParam)
End Function

' Initializes the hidden window and registers it with the Go core
Public Sub StartConfigWatcher(ByVal handle As LongPtr)
    If hProxyWnd <> 0 Then Exit Sub ' Guard: already running
    hUniLogHandle = handle
    
    ' 1. Create a "Message-Only" window (invisible, no graphical presence)
    hProxyWnd = CreateWindowExA(0, "Static", "UniLogProxy", 0, 0, 0, 0, 0, HWND_MESSAGE, 0, GetModuleHandleA(vbNullString), ByVal 0&)
    
    If hProxyWnd = 0 Then
        Debug.Print "!!! UniLog Error: Could not create Proxy Window"
        Exit Sub
    End If
    
    ' 2. Subclass the window to use our UniLog_WindowProc
    pOldWndProc = SetWindowLongPtrA(hProxyWnd, GWLP_WNDPROC, AddressOf UniLog_WindowProc)
    
    ' 3. Register our Window and the Specific Message ID with the Go DLL
    UniLog_RegisterVBAWindow handle, hProxyWnd, UNILOG_UPDATE_MSG
    
    Debug.Print "!!! UniLog: VBA Config Watcher Started (HWND: " & hProxyWnd & ")"
End Sub

' Safely stops the watcher and cleans up resources
Public Sub StopConfigWatcher()
    If hProxyWnd = 0 Then Exit Sub
    
    ' 1. Restore the original window procedure (Essential for stability)
    SetWindowLongPtrA hProxyWnd, GWLP_WNDPROC, pOldWndProc
    
    ' 2. Terminate the window
    DestroyWindow hProxyWnd
    
    hProxyWnd = 0
    pOldWndProc = 0
    Debug.Print "!!! UniLog: VBA Config Watcher Stopped."
End Sub

' -----------------------------------------------------------------------------
' WRAPPER FUNCTIONS
' -----------------------------------------------------------------------------

Public Function GetConfig(ByVal handle As LongPtr, ByVal section As String, ByVal key As String, Optional ByVal defaultVal As String = "") As String
    Dim ptr As LongPtr
    ptr = UniLog_Config_Get(handle, section, key)
    If ptr = 0 Then
        GetConfig = defaultVal
    Else
        GetConfig = PtrToString(ptr)
    End If
End Function

Public Sub SetConfig(ByVal handle As LongPtr, ByVal section As String, ByVal key As String, ByVal value As String)
    UniLog_Config_Set handle, section, key, value
End Sub

' -----------------------------------------------------------------------------
' DEMONSTRATION SUB
' -----------------------------------------------------------------------------

Public Sub TestUniversalLogger()
    Dim handle As LongPtr
    Dim dbIp As String
    
    ' 1. Initialize the logger
    handle = UniLog_Init("standalone", "Excel-Tool", "standard", Level_INFO)
    
    If handle = 0 Then
        MsgBox "Failed to initialize Universal Logger!", vbCritical
        Exit Sub
    End If
    
    ' 2. Log some messages
    UniLog_LogWithMetadata handle, Level_INFO, "Universal Logger initialized from Excel VBA", "UniversalLogger.bas", "77", "TestUniversalLogger", "VBA-Module"
    
    ' 3. Get configuration value (Using new standardized name)
    dbIp = GetConfig(handle, "database", "ip", "127.0.0.1")
    
    Debug.Print "Database IP from Config: " & dbIp
    UniLog_LogWithMetadata handle, Level_DEBUG, "Configured DB IP: " & dbIp, "UniversalLogger.bas", "83", "TestUniversalLogger", "VBA-Module"
    
    ' 4. Update configuration (Using new standardized name)
    SetConfig handle, "runtime", "last_run", Now()
    
    ' 5. Clean up
    UniLog_Close handle
    
    MsgBox "Logging complete! Check the logger output.", vbInformation
End Sub
