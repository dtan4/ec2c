package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestLaunchCommand_implement(t *testing.T) {
	var _ cli.Command = &LaunchCommand{}
}
