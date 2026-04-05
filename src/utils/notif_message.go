package utils

import (
	logger_models "github.com/Bastien-Antigravity/flexible-logger/src/models"
)

// -----------------------------------------------------------------------------

// NotifMessage mirrors the flexible-logger NotifMessage using a type alias.
// This allows notif-server to use the model without direct dependency on flexible-logger.
type NotifMessage = logger_models.NotifMessage
