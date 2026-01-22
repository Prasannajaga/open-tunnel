package controller

import (
	"log"

	"opentunnel/client/config"
	"opentunnel/client/constants"
	"opentunnel/client/service"
	"opentunnel/client/utils"
)

var TunnelHelp = utils.CommandHelp{
	Name:        "open tunnel",
	Description: "Opens a secure tunnel to the server",
	Usage:       "open tunnel [flags]",
	Options: []utils.OptionHelp{
		{
			Flag:        constants.FlagPort,
			ShortFlag:   constants.FlagPortChar,
			Description: "Port to bind the tunnel",
			Default:     "8080",
		},
		{
			Flag:        constants.FlagHelp,
			ShortFlag:   constants.FlagHelpChar,
			Description: "Show help",
		},
	},
	Examples: []string{
		"open tunnel",
		"open tunnel --port 9000",
	},
}

type TunnelController struct {
	cfg *config.Config
}

func NewTunnelController(cfg *config.Config) *TunnelController {
	return &TunnelController{cfg: cfg}
}

func (c *TunnelController) ShowHelp() {
	utils.PrintHelp(TunnelHelp)
}

func (c *TunnelController) Run() {
	svc := service.NewTunnelService(c.cfg)

	if err := svc.Connect(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	svc.Listen()
}
