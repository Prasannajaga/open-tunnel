package controller

import (
	"log"

	"opentunnel/server/config"
	"opentunnel/server/service"
)

type TunnelController struct {
	cfg *config.Config
}

func NewTunnelController(cfg *config.Config) *TunnelController {
	return &TunnelController{cfg: cfg}
}

func (c *TunnelController) Run() {
	svc := service.NewTunnelService(c.cfg)

	if err := svc.Start(); err != nil {
		log.Fatalf("Failed to start tunnel service: %v", err)
	}
}
