package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

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
		{
			Flag:        constants.FlagExport,
			Description: "Print export command for the token",
		},
	},
	Examples: []string{
		"open login",
		"open login --export",
	},
}

type AuthController struct {
	cfg         *config.Config
	ExportToken bool
}

func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{cfg: cfg}
}

func (c *AuthController) ShowHelp() {
	utils.PrintHelp(LoginHelp)
}

func (c *AuthController) Run() {
	// 1. Start CLI Auth Session
	startURL := fmt.Sprintf("http://%s:%s/api/cli-auth/start", c.cfg.ServerIP, c.cfg.HTTPPort)
	resp, err := http.Post(startURL, "application/json", nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	var startResult struct {
		SessionID string `json:"session_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&startResult); err != nil {
		log.Fatalf("Failed to parse server response: %v", err)
	}

	sessionID := startResult.SessionID
	loginURL := fmt.Sprintf("%s?session_id=%s", c.cfg.GetAuthURL(), sessionID)

	fmt.Println("Opening browser for authentication...")
	fmt.Printf("If the browser doesn't open, visit: %s\n", loginURL)

	if err := openBrowser(loginURL); err != nil {
		log.Printf("Could not open browser: %v", err)
	}

	fmt.Println("\nWaiting for login completion...")

	// 2. Poll for Token
	pollURL := fmt.Sprintf("http://%s:%s/api/cli-auth/poll/%s", c.cfg.ServerIP, c.cfg.HTTPPort, sessionID)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-timeout:
			log.Fatal("Login timed out after 5 minutes.")
		case <-ticker.C:
			resp, err := http.Get(pollURL)
			if err != nil {
				continue // Silently retry on network errors
			}

			if resp.StatusCode == http.StatusAccepted {
				resp.Body.Close()
				continue // Still pending
			}

			if resp.StatusCode == http.StatusOK {
				var result struct {
					Token string `json:"token"`
				}
				json.NewDecoder(resp.Body).Decode(&result)
				resp.Body.Close()

				// 3. Save Token
				if err := utils.SaveToken(result.Token); err != nil {
					log.Fatalf("Failed to save token: %v", err)
				}

				fmt.Println("✔ Login successful!")
				tokenPath, _ := utils.GetTokenPath()
				fmt.Printf("Token saved to %s\n", tokenPath)

				if c.ExportToken {
					fmt.Printf("\nTo use this token in your current shell, run:\n")
					fmt.Printf("  export OPENTUNNEL_TOKEN=%s\n", result.Token)
				}
				return
			}
			resp.Body.Close()
		}
	}
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
