package main

import (
	"dd_web/back/internal/app"
	"fmt"
	"os"
)

func main() {
	if err := app.RunOrder(); err != nil {
		fmt.Fprintf(os.Stderr, "order service exited: %v\n", err)
		os.Exit(1)
	}
}
