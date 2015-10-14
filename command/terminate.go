package command

import (
	"strings"
)

type TerminateCommand struct {
	Meta
}

func (c *TerminateCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *TerminateCommand) Synopsis() string {
	return "Terminate the specified EC2 instance"
}

func (c *TerminateCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
