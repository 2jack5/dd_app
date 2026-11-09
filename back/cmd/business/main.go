package main

import (
	"dd_web/back/internal/app"
	"fmt"
	"os"
)

func main() {
	if err := app.RunBusiness(); err != nil {
		fmt.Fprintf(os.Stderr, "business service exited: %v\n", err)
		os.Exit(1)
	}
}
