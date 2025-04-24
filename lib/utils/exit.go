package utils

import (
	"fmt"
	"os"
)

// ExitSuccess is a wrapper for exiting the execution of the program
// in case we want to add some custom cleanup or logging here
// without having to manually add this in every instance.
func ExitSuccess() {
	os.Exit(0)
}

func ExitError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
