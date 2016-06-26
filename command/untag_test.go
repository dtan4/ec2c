package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestUntagCommand_implement(t *testing.T) {
	var _ cli.Command = &TagCommand{}
}
