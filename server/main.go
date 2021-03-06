package main

import (
	"os"
	"strings"

	"github.com/SierraSoftworks/heimdall"
	"github.com/SierraSoftworks/heimdall/transports"
	log "github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"github.com/urfave/cli"
)

var config *Config

func main() {
	if envDSN := os.Getenv("SENTRY_DSN"); envDSN != "" {
		raven.SetDSN(envDSN)
	}

	raven.SetRelease(heimdall.Version)
	raven.SetEnvironment("server")

	app := cli.NewApp()
	app.Name = "Heimdall Server"
	app.Usage = "Run a Heimdall server instance"

	app.Author = "Benjamin Pannell"
	app.Email = "admin@sierrasoftworks.com"
	app.Copyright = "Sierra Softworks © 2016"
	app.Version = heimdall.Version

	app.Commands = cli.Commands{}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level",
			Usage: "DEBUG|INFO|WARN|ERROR",
			Value: "INFO",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "FILE",
			Value: "config.yaml",
		},
	}

	app.Before = func(c *cli.Context) error {
		log.WithFields(log.Fields{
			"log-level": c.GlobalString("log-level"),
		}).Info("Starting")

		logLevel := c.GlobalString("log-level")
		switch strings.ToUpper(logLevel) {
		case "DEBUG":
			log.SetLevel(log.DebugLevel)
		case "INFO":
			log.SetLevel(log.InfoLevel)
		case "WARN":
			log.SetLevel(log.WarnLevel)
		case "ERROR":
			log.SetLevel(log.ErrorLevel)
		default:
			log.SetLevel(log.InfoLevel)
		}

		cfg, err := ReadConfig(c.GlobalString("config"))
		if err != nil {
			return err
		}

		config = cfg

		return nil
	}

	app.Action = func(c *cli.Context) error {
		con := NewConsumer()

		for _, tc := range config.Transports {
			t, err := transports.GetTransport(tc.Driver, tc.URL)
			if err != nil {
				log.
					WithField("transport", tc).
					WithError(err).
					Error("Failed to connect to transport")
				continue
			}

			st, err := NewServerTransport(t)
			if err != nil {
				log.
					WithField("transport", tc).
					WithError(err).
					Error("Failed to configure server transport")
				continue
			}
			con.AddTransport(st)
		}

		go func() {
			con.Run()
		}()

		return ListenAndServe(config.Listen)
	}

	err := app.Run(os.Args)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		os.Exit(2)
	}

	os.Exit(0)
}
