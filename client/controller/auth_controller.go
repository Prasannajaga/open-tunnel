package controller

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"opentunnel/client/config"
	"opentunnel/client/constants"
	"opentunnel/client/utils"
)

var LoginHelp = utils.CommandHelp{
	Name:        "open login",
	Description: "Authenticate with the tunnel server",
	Usage:       "open login [flags]",
	Options: []utils.OptionHelp{
		{
			Flag:        constants.FlagHelp,
			ShortFlag:   constants.FlagHelpChar,
			Description: "Show help",
		},
	},
	Examples: []string{
		"open login",
	},
}

type AuthController struct {
	cfg *config.Config
}

func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{cfg: cfg}
}

func (c *AuthController) ShowHelp() {
	utils.PrintHelp(LoginHelp)
}

func (c *AuthController) Run() {
	loginURL := fmt.Sprintf("http://%s:%s/login", c.cfg.ServerIP, c.cfg.HTTPPort)

	fmt.Println("Opening browser for authentication...")
	fmt.Printf("If the browser doesn't open, visit: %s\n", loginURL)

	if err := openBrowser(loginURL); err != nil {
		log.Printf("Could not open browser: %v", err)
	}

	fmt.Println("\nAfter logging in, copy the token and run:")
	fmt.Println("  export OPENTUNNEL_TOKEN=<your-token>")
	fmt.Println("\nOr save it to ~/.opentunnel/token")
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
