Attribute VB_Name = "UniLog"
' -----------------------------------------------------------------------------
' UniLog - VBA Bridge for Microsoft Excel
' -----------------------------------------------------------------------------
' This module allows calling the Universal Logger (Go) shared library from VBA.
' Ensure that libunilog.dll is in the same directory as the workbook or in a
' directory listed in the system PATH.
' -----------------------------------------------------------------------------

Option Explicit

#If Win64 Then
    ' 64-bit Excel
    Private Declare PtrSafe Function UniLog_Init Lib "libunilog.dll" (ByVal configProfile As String, ByVal appName As String, ByVal loggerProfile As String, ByVal logLevel As Long) As LongPtr
    Private Declare PtrSafe Sub UniLog_Close Lib "libunilog.dll" (ByVal handle As LongPtr)
    Private Declare PtrSafe Function UniLog_Config_Get Lib "libunilog.dll" (ByVal handle As LongPtr, ByVal section As String, ByVal key As String) As LongPtr
    Private Declare PtrSafe Sub UniLog_Config_Set Lib "libunilog.dll" (ByVal handle As LongPtr, ByVal section As String, ByVal key As String, ByVal value As String)
    Private Declare PtrSafe Sub UniLog_LogWithMetadata Lib "libunilog.dll" (ByVal handle As LongPtr, ByVal level As LongLong, ByVal msg As String, ByVal file As String, ByVal line As String, ByVal functionName As String, ByVal moduleName As String)
    Private Declare PtrSafe Sub UniLog_SetLevel Lib "libunilog.dll" (ByVal handle As LongPtr, ByVal level As LongLong)
    
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
' WRAPPER FUNCTIONS (Standardized Names)
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
