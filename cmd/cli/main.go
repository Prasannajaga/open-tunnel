package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"opentunnel/client/config"
	"opentunnel/client/constants"
	"opentunnel/client/controller"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env")

	// database.TestConnection()
	// fmt.Println("hello")

	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case constants.CmdTunnel:
		handleTunnel(cmdArgs)
	case constants.CmdLogin:
		handleLogin(cmdArgs)
	case constants.FlagHelp, constants.FlagHelpChar:
		handleTunnel(args[0:])
	default:
		log.Fatalf(constants.ErrUnknownCommand, cmd)
	}
}

func handleTunnel(args []string) {
	cfg := config.NewConfig()

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch arg {
		case constants.FlagHelp, constants.FlagHelpChar:
			ctrl := controller.NewTunnelController(cfg)
			ctrl.ShowHelp()
			return

		case constants.FlagPort, constants.FlagPortChar:
			if i+1 >= len(args) {
				log.Fatal("missing value for --port")
			}
			port, err := strconv.Atoi(args[i+1])
			if err != nil {
				log.Fatalf(constants.ErrInvalidPort, args[i+1])
			}
			cfg.LocalPort = port
			i++
		}
	}

	ctrl := controller.NewTunnelController(cfg)
	ctrl.Run()
}

func handleLogin(args []string) {
	cfg := config.NewConfig()
	exportToken := false

	for _, arg := range args {
		switch arg {
		case constants.FlagHelp, constants.FlagHelpChar:
			ctrl := controller.NewAuthController(cfg)
			ctrl.ShowHelp()
			return
		case constants.FlagExport:
			exportToken = true
		}
	}

	ctrl := controller.NewAuthController(cfg)
	ctrl.ExportToken = exportToken
	ctrl.Run()
}

func printUsage() {
	fmt.Println("Usage: open <command> [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  tunnel    Open a secure tunnel to the server")
	fmt.Println("  login     Authenticate with the tunnel server")
	fmt.Println()
	fmt.Println("Use 'open <command> --help' for more information about a command.")
}
