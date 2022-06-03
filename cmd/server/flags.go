package main

import (
	"github.com/urfave/cli/v2"
)

const (
	hostFlag        = "host"
	portFlag        = "port"
	defaultHTTPHost = "0.0.0.0"
	defaultHTTPPort = "9002"
)

type WithFlagsFunc func() []cli.Flag

func CLIFlags(options ...WithFlagsFunc) []cli.Flag {
	flags := []cli.Flag{}

	for _, group := range options {
		flags = append(flags, group()...)
	}

	return flags
}

func parseFlags(cfg *appConfig) cli.BeforeFunc {
	return func(ctx *cli.Context) error {
		cfg.Host = ctx.String(hostFlag)
		cfg.Port = ctx.String(portFlag)

		return nil
	}
}

func WithHTTPServerFlags() WithFlagsFunc {
	return func() []cli.Flag {
		return []cli.Flag{
			&cli.StringFlag{
				Name:  hostFlag,
				Value: defaultHTTPHost,
				Usage: "HTTP listening host",
			},
			&cli.StringFlag{
				Name:  portFlag,
				Value: defaultHTTPPort,
				Usage: "HTTP listening port",
			},
		}
	}
}
