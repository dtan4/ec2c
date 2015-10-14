package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestTerminateCommand_implement(t *testing.T) {
	var _ cli.Command = &TerminateCommand{}
}
