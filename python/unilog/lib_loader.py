#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from ctypes import CDLL as ctypeCDLL, CFUNCTYPE as ctypeCFUNCTYPE, c_char_p as ctypeC_char_p, \
                   c_int as ctypeC_int, c_size_t as ctypeC_size_t, c_longlong as ctypeC_longlong
from ctypes.util import find_library as ctypeUtilFindLibrary
from pathlib import Path as pathlibPath


##########################################################################
# Loader logic

# Discovery function to find the shared library across development and production environments
def _load_lib():
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
        for ext in [".dylib", ".dll"]:
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

if lib:
    # Initialization & Lifecycle
    lib.UniLog_Init.argtypes = [ctypeC_char_p, ctypeC_char_p, ctypeC_char_p, ctypeC_int]
    lib.UniLog_Init.restype = ctypeC_size_t
    
    lib.UniLog_Close.argtypes = [ctypeC_size_t]

    # Logging Interface
    lib.UniLog_LogWithMetadata.argtypes = [
        ctypeC_size_t, ctypeC_longlong, ctypeC_char_p, 
        ctypeC_char_p, ctypeC_char_p, ctypeC_char_p, ctypeC_char_p
    ]
    lib.UniLog_SetLevel.argtypes = [ctypeC_size_t, ctypeC_longlong]

    # Configuration Interface
    lib.UniLog_Config_Get.argtypes = [ctypeC_size_t, ctypeC_char_p, ctypeC_char_p]
    lib.UniLog_Config_Get.restype = ctypeC_char_p
    
    lib.UniLog_Config_Set.argtypes = [ctypeC_size_t, ctypeC_char_p, ctypeC_char_p, ctypeC_char_p]

    # Shared Bridge Callbacks
    CALLBACK_TYPE = ctypeCFUNCTYPE(None, ctypeC_char_p)
    lib.UniLog_OnMemConfUpdate.argtypes = [ctypeC_size_t, CALLBACK_TYPE]
