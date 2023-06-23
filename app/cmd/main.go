package main

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v8"
	"github.com/urfave/cli/v2"

	"github.com/abtergo/abtergo/app/config"
	"github.com/abtergo/abtergo/app/http"
	"github.com/abtergo/abtergo/libs/cleaner"
)

func main() {
	logger := createLogger()
	cache := createCache()

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

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
			return http.NewServer(cfg, logger, cleaner.New(logger), cache).
				SetupMiddleware(cCtx).
				SetupHandlers().
				Start(cCtx)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
