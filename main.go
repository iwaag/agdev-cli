package main

import (
	"context"
	"fmt"
	"os"

	"agdev/cmd"
	"agdev/internal/app"
)

func main() {
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(app.ExitCode(err))
	}
}
