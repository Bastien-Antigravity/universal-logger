package interfaces

import (
	distributed_config "github.com/Bastien-Antigravity/distributed-config"
	"github.com/Bastien-Antigravity/flexible-logger/src/interfaces"
)

// IFacade is the central orchestrator combining configuration and logging.
type IFacade interface {
	interfaces.Logger
	GetConfig() *distributed_config.Config
}
