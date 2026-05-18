#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from ctypes import CDLL as ctypeCDLL, CFUNCTYPE as ctypeCFUNCTYPE, c_char_p as ctypeC_char_p, \
                   c_int as ctypeC_int, c_size_t as ctypeC_size_t, c_longlong as ctypeC_longlong, \
                   c_void_p as ctypeC_void_p
from ctypes.util import find_library as ctypeUtilFindLibrary
from pathlib import Path as pathlibPath
from typing import Any


##########################################################################
# Loader logic

# Discovery function to find the shared library across development and production environments
def _load_lib() -> Any:
    lib_name = "libunilog"
    
    # 1. Check local package directory (for distributed wheels)
    lib_path = pathlibPath(__file__).parent / f"{lib_name}.so"
    
    # 2. Check root 'libunilog' directory (for development)
    root_lib = pathlibPath(__file__).parent.parent.parent / "libunilog" / f"{lib_name}.so"

    found = False
    if lib_path.exists():
        found = True
    elif root_lib.exists():
        lib_path = root_lib
        found = True
        
    # Resolve platform-specific extensions if the generic .so is not found
    if not found:
        # Prioritize extension based on current platform
        import sys
        platforms = [".dll", ".dylib"] if sys.platform == "win32" else [".dylib", ".so"]
        for ext in platforms:
            p = lib_path.with_suffix(ext)
            if p.exists():
                lib_path = p
                found = True
                break
            p_root = root_lib.with_suffix(ext)
            if p_root.exists():
                lib_path = p_root
                found = True
                break
    
    if not found:
        # Fallback to system path lookup using standard OS utilities
        res = ctypeUtilFindLibrary(lib_name)
        print(f"!!! libunilog: System fallback -> {res}")
        return ctypeCDLL(res or lib_name)
    
    print(f"!!! libunilog: Loading from -> {lib_path}")
    return ctypeCDLL(str(lib_path))


##########################################################################
# Library loading

try:
    lib = _load_lib()
except Exception:
    # Silent failure during initialization to avoid blocking installation or builds
    lib = None


##########################################################################
# FFI Declarations

# Shared Bridge Callbacks (defined globally to prevent ImportErrors)
CALLBACK_TYPE = ctypeCFUNCTYPE(None, ctypeC_char_p)

if lib:
    # Initialization & Lifecycle
    lib.UniLog_Init.argtypes = [
        ctypeC_char_p, ctypeC_char_p, ctypeC_char_p, 
        ctypeC_int, ctypeC_int, ctypeC_size_t
    ]
    lib.UniLog_Init.restype = ctypeC_size_t
    lib.UniLog_Close.argtypes = [ctypeC_size_t]

    # Logging Interface
    lib.UniLog_LogWithMetadata.argtypes = [
        ctypeC_size_t, ctypeC_longlong, ctypeC_char_p, 
        ctypeC_char_p, ctypeC_char_p, ctypeC_char_p, ctypeC_char_p
    ]
    lib.UniLog_SetLevel.argtypes = [ctypeC_size_t, ctypeC_longlong]
    lib.UniLog_GetLevel.argtypes = [ctypeC_size_t]
    lib.UniLog_GetLevel.restype = ctypeC_int

    # Metadata Interface
    lib.UniLog_AddMetadata.argtypes = [ctypeC_size_t, ctypeC_char_p, ctypeC_char_p]
    lib.UniLog_SetMetadata.argtypes = [ctypeC_size_t, ctypeC_char_p]

    # Configuration Interface (Native UniLog)
    lib.DistConf_FreeString.argtypes = [ctypeC_void_p]
    lib.UniLog_Config_Get.argtypes = [ctypeC_size_t, ctypeC_char_p, ctypeC_char_p]
    lib.UniLog_Config_Get.restype = ctypeC_void_p
    lib.UniLog_Config_Set.argtypes = [ctypeC_size_t, ctypeC_char_p, ctypeC_char_p, ctypeC_char_p]
    lib.UniLog_OnConfigUpdate.argtypes = [ctypeC_size_t, CALLBACK_TYPE]
    lib.UniLog_RegisterNotifCallback.argtypes = [ctypeC_size_t, CALLBACK_TYPE]

    # DISTCONF COMPATIBILITY LAYER (Microservice Toolbox Support)
    try:
        lib.DistConf_New.argtypes = [ctypeC_char_p]
        lib.DistConf_New.restype = ctypeC_size_t
        lib.DistConf_Get.argtypes = [ctypeC_size_t, ctypeC_char_p, ctypeC_char_p]
        lib.DistConf_Get.restype = ctypeC_void_p
        lib.DistConf_Set.argtypes = [ctypeC_size_t, ctypeC_char_p, ctypeC_char_p, ctypeC_char_p]
        lib.DistConf_GetFullConfig.argtypes = [ctypeC_size_t]
        lib.DistConf_GetFullConfig.restype = ctypeC_void_p
        lib.DistConf_Close.argtypes = [ctypeC_size_t]
    except Exception:
        pass # Compatibility layer might be missing if lib is old
