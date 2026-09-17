#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
=============================================================================
ESSENTIAL PROCESS:
Package initialization for the unilog Python library, exposing primary classes
and enumerations for external consumer import.

DATA FLOW:
1. Input: Submodule symbols from facade, models, and listeners.
2. Logic: Exports public API surface via __all__.
3. Output: Unified import namespace (UniLog, LogLevel).

KEY PARAMETERS:
- None
=============================================================================
"""


from .facade import UniLog
from .models import LogLevel

__all__ = ['UniLog', 'LogLevel']
