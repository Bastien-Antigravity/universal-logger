package factory

import (
	"github.com/Bastien-Antigravity/distconf-flexlog/src/facade"
	"github.com/Bastien-Antigravity/distconf-flexlog/src/interfaces"
	"github.com/Bastien-Antigravity/distconf-flexlog/src/models"
)

// NewFacade creates a new fully configured DistconfFlexlogFacade.
func NewFacade(p models.MFacadeParams) interfaces.IFacade {
	return facade.NewDistconfFlexlogFacade(p)
}
