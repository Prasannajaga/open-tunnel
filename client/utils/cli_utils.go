package utils

import (
	"fmt"
	"strings"
)

type CommandHelp struct {
	Name        string
	Description string
	Usage       string
	Options     []OptionHelp
	Examples    []string
}

type OptionHelp struct {
	Flag        string
	ShortFlag   string
	Description string
	Default     string
}

func PrintHelp(cmd CommandHelp) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Command:\n  %s\n\n", cmd.Name))
	sb.WriteString(fmt.Sprintf("Description:\n  %s\n\n", cmd.Description))
	sb.WriteString(fmt.Sprintf("Usage:\n  %s\n\n", cmd.Usage))

	if len(cmd.Options) > 0 {
		sb.WriteString("Options:\n")
		for _, opt := range cmd.Options {
			flagStr := opt.Flag
			if opt.ShortFlag != "" {
				flagStr += opt.ShortFlag
			}
			desc := opt.Description
			if opt.Default != "" {
				desc += fmt.Sprintf(" (default: %s)", opt.Default)
			}
			sb.WriteString(fmt.Sprintf("  %-14s %s\n", flagStr, desc))
		}
		sb.WriteString("\n")
	}

	if len(cmd.Examples) > 0 {
		sb.WriteString("Examples:\n")
		for _, ex := range cmd.Examples {
			sb.WriteString(fmt.Sprintf("  %s\n", ex))
		}
	}

	fmt.Print(sb.String())
}
