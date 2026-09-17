package utils

// =============================================================================
// ESSENTIAL PROCESS: Type aliases mapping notification messages to interface definitions.
//
// DATA FLOW:
//   1. Mirrors interfaces.NotifMessage in utils package for consumer convenience.
//
// KEY PARAMETERS:
//   - NotifMessage: Alert message type alias.
// =============================================================================


import (
	"github.com/Bastien-Antigravity/universal-logger/src/interfaces"
)

// -----------------------------------------------------------------------------

// NotifMessage mirrors the universal-logger NotifMessage using a type alias.
type NotifMessage = interfaces.NotifMessage
