#!/usr/bin/env python3
# -*- coding: utf-8 -*-

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
