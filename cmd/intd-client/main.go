package main

import (
	"context"
	"github.com/aaronland/go-artisanal-integers/application"
	"log"
	"os"
)

func main() {

	flags := application.NewClientApplicationFlags()

	ctx := context.Background()
	app, err := application.NewClientApplication(ctx, "client://")

	if err != nil {
		log.Fatal(err)
	}

	err = app.Run(ctx, flags)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
