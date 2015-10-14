package command

import (
	"strings"
)

type LaunchCommand struct {
	Meta
}

func (c *LaunchCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *LaunchCommand) Synopsis() string {
	return ""
}

func (c *LaunchCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
