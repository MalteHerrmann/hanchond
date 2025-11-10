package utils

import (
	"fmt"
	"strings"
)

// Log handles the printing of messages to the console.
// At present, this just wraps fmt.Printf and fmt.Println.
// It may be extended to add more features in the future / other types of logging.
func Log(message string, args ...any) {
	if len(args) > 0 {
		if !strings.HasSuffix(message, "\n") {
			message += "\n"
		}

		fmt.Printf(message, args...)
	} else {
		fmt.Println(message)
	}
}
