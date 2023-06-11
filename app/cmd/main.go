package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	cleaner2 "github.com/abtergo/abtergo/libs/cleaner"
)

func main() {
	logger := createLogger()

	app := &cli.App{
		Name:  "server",
		Usage: "start an HTTP server",
		Action: func(cCtx *cli.Context) error {
			return NewHTTPServer(logger, cleaner2.New(logger)).SetupHandlers(logger).Start(cCtx)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
