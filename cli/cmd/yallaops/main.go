// Command yallaops is the YallaOps CLI.
package main

import (
	"fmt"
	"os"

	"github.com/yallaops/yallaops/cli/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
