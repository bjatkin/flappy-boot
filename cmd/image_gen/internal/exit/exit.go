package exit

import (
	"fmt"
	"os"
)

// Error codes for use in this package
const (
	InvalidArguments = iota + 1
	InvalidConfig
	InvalidPalette
	InvalidTileSet
	InvalidTileMap
	FileWriteFailed
)

var (
	exitCode int
	exitErr  error
)

// Final should be run as the first defer in the main function. It prints the currently set
// exit error and then exits with the correct error code. If the exit error is nil Final is a no-op
func Final() {
	if exitErr != nil {
		fmt.Println(exitErr.Error())
		os.Exit(exitCode)
	}
}

// Error sets the current exit code and exit error. It should be called from the main function and
// you should imediatly return afterwards
func Error(code int, err error) {
	exitCode = code
	exitErr = err
}
