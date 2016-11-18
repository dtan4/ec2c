package msg

import (
	"fmt"
	"os"
)

// Errorf shows error message with the given format
func Errorf(format string, err error) {
	fmt.Fprintf(os.Stderr, format, err)
}
