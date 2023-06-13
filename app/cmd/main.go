package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/abtergo/abtergo/app/http"
	"github.com/abtergo/abtergo/libs/cleaner"
)

func main() {
	logger := createLogger()

	app := &cli.App{
		Name:  "server",
		Usage: "start an HTTP server",
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   8080,
				Usage:   "port to listen on",
			},
		},
		Action: func(cCtx *cli.Context) error {
			return http.NewServer(logger, cleaner.New(logger)).
				SetupMiddleware(cCtx).
				SetupHandlers().
				Start(cCtx)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
