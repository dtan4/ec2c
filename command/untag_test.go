package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestTagCommand_implement(t *testing.T) {
	var _ cli.Command = &TagCommand{}
}
