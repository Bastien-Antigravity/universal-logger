#!/usr/bin/env python
# coding:utf-8
"""
=============================================================================
ESSENTIAL PROCESS:
Automated cross-language parity audit script verifying that all exported CGO
symbols from libunilog.h are implemented across Python, Rust, C++, and VBA.

DATA FLOW:
1. Input: libunilog.h C header declarations and language wrapper source files.
2. Logic: Parses exported UniLog_* functions and verifies matching bindings in each language.
3. Output: Comprehensive parity report table across Python, Rust, C++, and VBA.

KEY PARAMETERS:
- None
=============================================================================
"""

from os.path import abspath as osPathAbspath, dirname as osPathDirname, exists as osPathExists, join as osPathJoin
from re import IGNORECASE as reIGNORECASE, findall as reFindall, search as reSearch
from sys import exit as sysExit
from typing import Dict as typeDict, List as typeList

# -----------------------------------------------------------------------------------------------
# CONFIGURATION & PATHS
# -----------------------------------------------------------------------------------------------

ROOT_DIR = osPathDirname(osPathDirname(osPathAbspath(__file__)))

PRIMARY_HEADER = osPathJoin(ROOT_DIR, "libunilog", "libunilog.h")
FALLBACK_HEADER = osPathJoin(ROOT_DIR, "unilog", "libunilog", "libunilog.h")
SOURCE_OF_TRUTH = PRIMARY_HEADER if osPathExists(PRIMARY_HEADER) else FALLBACK_HEADER

FACADES = {
    "Python": {
        "file": osPathJoin(ROOT_DIR, "unilog", "python", "unilog", "facade.py"),
        "pattern": r"def\s+([a-zA-Z0-9_]+)\("
    },
    "Rust": {
        "file": osPathJoin(ROOT_DIR, "unilog", "rust", "src", "lib.rs"),
        "pattern": r"pub\s+fn\s+([a-z0-9_]+)\("
    },
    "C++": {
        "file": osPathJoin(ROOT_DIR, "unilog", "cpp", "UniversalLogger.hpp"),
        "pattern": r"([a-zA-Z0-9_]+)\("
    },
    "VBA": {
        "file": osPathJoin(ROOT_DIR, "unilog", "vba", "UniversalLogger.bas"),
        "pattern": r"(?:Function|Sub)\s+([a-zA-Z0-9_]+)\("
    }
}

EXPECTED_LEVELS = [
    "DEBUG", "STREAM", "INFO", "LOGON", "LOGOUT", "TRADE", "SCHEDULE", "REPORT", "WARNING", "ERROR", "CRITICAL"
]

EXPORT_MAPPINGS = {
    "UniLog_Init": {"Python": "__init__", "Rust": "new", "C++": "UniLog", "VBA": "UniLog_Init"},
    "UniLog_Close": {"Python": "close", "Rust": "drop", "C++": "~UniLog", "VBA": "UniLog_Close"},
    "UniLog_LogWithMetadata": {"Python": "_dispatch_log_to_cgo", "Rust": "log_with_metadata", "C++": "log_with_metadata", "VBA": "UniLog_LogWithMetadata"},
    "UniLog_SetLevel": {"Python": "set_level", "Rust": "set_level", "C++": "set_level", "VBA": "UniLog_SetLevel"},
    "UniLog_GetLevel": {"Python": "get_level", "Rust": "get_level", "C++": "get_level", "VBA": "GetLevel"},
    "UniLog_AddMetadata": {"Python": "add_metadata", "Rust": "add_metadata", "C++": "add_metadata", "VBA": "AddMetadata"},
    "UniLog_SetMetadata": {"Python": "set_metadata", "Rust": "set_metadata", "C++": "set_metadata", "VBA": "UniLog_SetMetadata"},
    "UniLog_Config_Get": {"Python": "get_config", "Rust": "get_config", "C++": "get_config", "VBA": "GetConfig"},
    "UniLog_Config_Get_Safe": {"Python": "get_config", "Rust": "get_config", "C++": "get_config", "VBA": "GetConfig"},
    "UniLog_Config_Set": {"Python": "set_config", "Rust": "set_config", "C++": "set_config", "VBA": "SetConfig"},
    "UniLog_OnConfigUpdate": {"Python": "on_config_update", "Rust": "on_config_update", "C++": "on_config_update", "VBA": "StartConfigWatcher"},
    "UniLog_RegisterNotifCallback": {"Python": "on_notification", "Rust": "on_notification", "C++": "set_notification_callback", "VBA": "UniLog_RegisterNotifCallback"},
    "UniLog_RegisterVBAWindow": {"Python": "__init__", "Rust": "new", "C++": "UniLog", "VBA": "StartConfigWatcher"},
    "UniLog_Flush": {"Python": "close", "Rust": "drop", "C++": "~UniLog", "VBA": "UniLog_Close"},
}

# -----------------------------------------------------------------------------------------------


def get_exported_functions(header_path: str) -> typeList[str]:
    """Extract UniLog_* exports from the C header."""
    if not osPathExists(header_path):
        print(f"Error: Header not found at {header_path}")
        return []

    with open(header_path, 'r') as f:
        content = f.read()

    pattern = r"extern\s+(?:void|char\*|GoUintptr|GoInt)\s+(UniLog_[a-zA-Z0-9_]+)\("
    matches = reFindall(pattern, content)
    return sorted(list(set(matches)))

# -----------------------------------------------------------------------------------------------


def check_facade_parity(facade_name: str, facade_info: dict, exported_funcs: typeList[str]) -> typeDict[str, bool]:
    """Check which exported functions are implemented in the given facade."""
    if not osPathExists(facade_info["file"]):
        return {f: False for f in exported_funcs}

    with open(facade_info["file"], 'r') as f:
        content = f.read()

    implemented = reFindall(facade_info["pattern"], content)
    implemented_set = {i.lower() for i in implemented}

    results = {}
    for func in exported_funcs:
        if facade_name in ("C++", "VBA") and func == "UniLog_Close":
            results[func] = ("~unilog" in content.lower()) if facade_name == "C++" else ("unilog_close" in content.lower())
            continue
        mapped = EXPORT_MAPPINGS.get(func, {}).get(facade_name, "")
        normalized = func.replace("UniLog_", "").lower()
        results[func] = (mapped != "" and mapped.lower() in implemented_set) or (normalized in implemented_set) or (func.lower() in content.lower())

    return results

# -----------------------------------------------------------------------------------------------


def check_level_parity(facade_info: dict) -> typeDict[str, bool]:
    """Check if all 11 log levels are defined in the facade."""
    if not osPathExists(facade_info["file"]):
        return {lvl: False for lvl in EXPECTED_LEVELS}

    with open(facade_info["file"], 'r') as f:
        content = f.read()

    results = {}
    for level in EXPECTED_LEVELS:
        patterns = [
            rf"\b{level}\b",
            rf"\bLevel_{level.capitalize()}\b",
            rf"\bLevel_{level}\b",
            rf"\b{level.lower()}\b"
        ]
        results[level] = any(reSearch(p, content, reIGNORECASE) for p in patterns)

    return results

# -----------------------------------------------------------------------------------------------


def main() -> None:
    print("\n" + "=" * 80)
    print(" BASTIEN-ANTIGRAVITY: POLYGLOT PARITY AUDIT")
    print("=" * 80)

    exports = get_exported_functions(SOURCE_OF_TRUTH)
    if not exports:
        sysExit(1)

    print(f"Detected {len(exports)} exported UniLog_* functions in {SOURCE_OF_TRUTH}\n")

    facade_results = {name: check_facade_parity(name, info, exports) for name, info in FACADES.items()}
    level_results = {name: check_level_parity(info) for name, info in FACADES.items()}

    header = f"{'Function':<30} | {'Py':<3} | {'Rs':<3} | {'C++':<3} | {'VBA':<3}"
    print(header)
    print("-" * len(header))

    for func in exports:
        row = f"{func:<30} | "
        for name in FACADES.keys():
            status = "✅" if facade_results[name][func] else "❌"
            row += f"{status:<3} | "
        print(row)

    print("\n" + "=" * 80)
    print(" LOG LEVEL PARITY")
    print("=" * 80)
    header_levels = f"{'Level':<30} | {'Py':<3} | {'Rs':<3} | {'C++':<3} | {'VBA':<3}"
    print(header_levels)
    print("-" * len(header_levels))

    for level in EXPECTED_LEVELS:
        row = f"{level:<30} | "
        for name in FACADES.keys():
            status = "✅" if level_results[name][level] else "❌"
            row += f"{status:<3} | "
        print(row)

    print("\nAudit Complete.\n")


if __name__ == "__main__":
    main()
