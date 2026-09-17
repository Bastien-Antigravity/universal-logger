#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
=============================================================================
ESSENTIAL PROCESS:
Data models, enumerations, and type definitions representing log levels,
metadata containers, and FFI callback signatures for UniLog.

DATA FLOW:
1. Input: Domain log level definitions and C types.
2. Logic: Maps Python Enum values to Go core integer constants.
3. Output: Typed models consumed by the UniLog facade and test runners.

KEY PARAMETERS:
- LogLevel: Enumeration of canonical log levels (Debug=1 through Critical=11).
=============================================================================
"""


from enum import IntEnum


##########################################################################
# Log Levels

class LogLevel(IntEnum):
    """
    LogLevel constants matching the Go logger_models.Level.
    """
    DEBUG = 1
    STREAM = 2
    INFO = 3
    LOGON = 4
    LOGOUT = 5
    TRADE = 6
    SCHEDULE = 7
    REPORT = 8
    WARNING = 9
    ERROR = 10
    CRITICAL = 11


    ##########################################################################
    # Parsers

    # Helper to convert string-based levels from config to enum
    @classmethod
    def from_str(cls, s: str) -> 'LogLevel':
        return getattr(cls, s.upper(), cls.INFO)
