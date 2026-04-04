#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import ctypes
import ctypes.util
from pathlib import Path

def _load_lib():
    lib_name = "libunilog"
    # 1. Check local package directory (for distributed wheels)
    lib_path = Path(__file__).parent / f"{lib_name}.so"
    
    # 2. Check root 'libunilog' directory (for development)
    root_lib = Path(__file__).parent.parent / "libunilog" / f"{lib_name}.so"

    found = False
    if lib_path.exists():
        found = True
    elif root_lib.exists():
        lib_path = root_lib
        found = True
        
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
        # Fallback to system path if not found in package directory
        res = ctypes.util.find_library(lib_name)
        return ctypes.CDLL(res or lib_name)
    
    return ctypes.CDLL(str(lib_path))

try:
    lib = _load_lib()
except Exception:
    # During build or if library is missing, we don't want to crash on import
    lib = None

if lib:
    # --- Initialization & Lifecycle ---
    lib.UniLog_Init.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
    lib.UniLog_Init.restype = ctypes.c_size_t
    
    lib.UniLog_Close.argtypes = [ctypes.c_size_t]

    # --- Logging ---
    lib.UniLog_LogWithMetadata.argtypes = [
        ctypes.c_size_t, ctypes.c_longlong, ctypes.c_char_p, 
        ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p
    ]
    lib.UniLog_SetLevel.argtypes = [ctypes.c_size_t, ctypes.c_longlong]

    # --- Configuration ---
    lib.UniLog_Config_Get.argtypes = [ctypes.c_size_t, ctypes.c_char_p, ctypes.c_char_p]
    lib.UniLog_Config_Get.restype = ctypes.c_char_p
    
    lib.UniLog_Config_Set.argtypes = [ctypes.c_size_t, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]

    # --- Callbacks ---
    CALLBACK_TYPE = ctypes.CFUNCTYPE(None, ctypes.c_char_p)
    lib.UniLog_OnMemConfUpdate.argtypes = [ctypes.c_size_t, CALLBACK_TYPE]
