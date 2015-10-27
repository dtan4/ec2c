package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestSpotRequestCommand_implement(t *testing.T) {
	var _ cli.Command = &SpotRequestCommand{}
}
