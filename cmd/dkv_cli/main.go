package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var (
	restyClient = resty.New()
	logger      = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()
)

func main() {
	app := cli.App{
		Name:    "dkv_cli",
		Version: "1.0.0",
		Usage:   "CLI to operate with distributed KV storage",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "endpoint",
				Usage: "Endpoint of DKV to operate with",
				Value: "localhost:8001",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "Get KV for specified key",
				Action:  get,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Aliases:  []string{"k"},
						Required: true,
					},
				},
			},
			{
				Name:    "put",
				Aliases: []string{"p"},
				Usage:   "Put/update KV",
				Action:  put,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Aliases:  []string{"k"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "value",
						Aliases:  []string{"val", "v"},
						Required: true,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal().
			Err(fmt.Errorf("running CLI: %w", err)).
			Send()
	}
}
