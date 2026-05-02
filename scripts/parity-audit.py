#!/usr/bin/env python
# coding:utf-8

from os.path import join as osPathJoin, dirname as osPathDirname, abspath as osPathAbspath, exists as osPathExists
from sys import exit as sysExit
from re import findall as reFindall, search as reSearch, IGNORECASE as reIGNORECASE
from typing import List, Dict, Set

# -----------------------------------------------------------------------------------------------
# CONFIGURATION & PATHS
# -----------------------------------------------------------------------------------------------

ROOT_DIR = osPathDirname(osPathDirname(osPathAbspath(__file__)))

SOURCE_OF_TRUTH = osPathJoin(ROOT_DIR, "libunilog", "libunilog.h")

FACADES = {
    "Python": {
        "file": osPathJoin(ROOT_DIR, "python", "unilog", "facade.py"),
        "type": "def",
        "pattern": r"def\s+([a-zA-Z0-9_]+)\("
    },
    "Rust": {
        "file": osPathJoin(ROOT_DIR, "rust", "src", "lib.rs"),
        "type": "fn",
        "pattern": r"pub\s+fn\s+([a-z0-9_]+)\("
    },
    "C++": {
        "file": osPathJoin(ROOT_DIR, "cpp", "UniversalLogger.hpp"),
        "type": "method",
        "pattern": r"([a-zA-Z0-9_]+)\("
    },
    "VBA": {
        "file": osPathJoin(ROOT_DIR, "vba", "UniversalLogger.bas"),
        "type": "function",
        "pattern": r"(?:Function|Sub)\s+([a-zA-Z0-9_]+)\("
    }
}

EXPECTED_LEVELS = [
    "DEBUG", "STREAM", "INFO", "LOGON", "LOGOUT", "TRADE", "SCHEDULE", "REPORT", "WARNING", "ERROR", "CRITICAL"
]

# -----------------------------------------------------------------------------------------------
# PARSING LOGIC
# -----------------------------------------------------------------------------------------------

def get_exported_functions(header_path: str) -> List[str]:
    """Extract UniLog_* and DistConf_* exports from the C header."""
    if not osPathExists(header_path):
        print(f"Error: Header not found at {header_path}")
        return []
    
    with open(header_path, 'r') as f:
        content = f.read()
    
    # Match extern void UniLog_... or extern char* DistConf_...
    pattern = r"extern\s+(?:void|char\*|GoUintptr|GoInt)\s+([a-zA-Z0-9_]+)\("
    matches = reFindall(pattern, content)
    return sorted(list(set(matches)))

def check_facade_parity(facade_name: str, facade_info: dict, exported_funcs: List[str]) -> Dict[str, bool]:
    """Check which exported functions are implemented in the given facade."""
    if not osPathExists(facade_info["file"]):
        return {f: False for f in exported_funcs}
    
    with open(facade_info["file"], 'r') as f:
        content = f.read()
    
    implemented = reFindall(facade_info["pattern"], content)
    implemented_set = {i.lower() for i in implemented}
    
    results = {}
    for func in exported_funcs:
        # Normalize Go function name (strip UniLog_ and DistConf_ and make snake_case)
        normalized = func.replace("UniLog_", "").replace("DistConf_", "").lower()
        
        # Special mappings if needed
        mappings = {
            "onliveconfupdate": "on_config_update", # Rust/Python convention
            "onconfigupdate": "on_config_update",
            "logwithmetadata": "log",
            "registernotifcallback": "on_notification",
            "registervbawindow": "register_vba_window",
            "init": "new" if facade_name == "Rust" else "init"
        }
        
        target = mappings.get(normalized, normalized)
        
        # Check for both normalized and target
        results[func] = (normalized in implemented_set or target in implemented_set)
        
    return results

def check_level_parity(facade_name: str, facade_info: dict) -> Dict[str, bool]:
    """Check if all 11 log levels are defined in the facade."""
    if not osPathExists(facade_info["file"]):
        return {l: False for l in EXPECTED_LEVELS}
    
    with open(facade_info["file"], 'r') as f:
        content = f.read()
    
    results = {}
    for level in EXPECTED_LEVELS:
        # Check for LEVEL_INFO, Level_Info, Info, etc.
        patterns = [
            rf"\b{level}\b",
            rf"\bLevel_{level.capitalize()}\b",
            rf"\bLevel_{level}\b",
            rf"\b{level.lower()}\b"
        ]
        found = any(reSearch(p, content, reIGNORECASE) for p in patterns)
        results[level] = found
        
    return results

# -----------------------------------------------------------------------------------------------
# EXECUTION
# -----------------------------------------------------------------------------------------------

def main():
    print("\n" + "="*80)
    print(" BASTIEN-ANTIGRAVITY: POLYGLOT PARITY AUDIT")
    print("="*80)
    
    exports = get_exported_functions(SOURCE_OF_TRUTH)
    if not exports:
        sysExit(1)
        
    print(f"Detected {len(exports)} exported functions in {SOURCE_OF_TRUTH}\n")
    
    # Audit Functions
    facade_results = {}
    for name, info in FACADES.items():
        facade_results[name] = check_facade_parity(name, info, exports)
        
    # Audit Levels
    level_results = {}
    for name, info in FACADES.items():
        level_results[name] = check_level_parity(name, info)

    # Print Table
    header = f"{'Function':<30} | {'Py':<3} | {'Rs':<3} | {'C++':<3} | {'VBA':<3}"
    print(header)
    print("-" * len(header))
    
    for func in exports:
        row = f"{func:<30} | "
        for name in FACADES.keys():
            status = "✅" if facade_results[name][func] else "❌"
            row += f"{status:<3} | "
        print(row)
        
    print("\n" + "="*80)
    print(" LOG LEVEL PARITY")
    print("="*80)
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
