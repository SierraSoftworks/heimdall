package main

import (
	"strings"

	"github.com/SierraSoftworks/heimdall/internal/app/heimdall"
	"github.com/SierraSoftworks/heimdall/internal/pkg/config"
	"github.com/SierraSoftworks/sentry-go"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var version = "1.0.0-dev"

func main() {
	raven := sentry.NewClient(
		sentry.Release(version),
	)

	app := cli.NewApp()
	app.Name = "Heimdall"
	app.Usage = "A robust, scalable and straightforward monitoring platform for binary state checks."

	app.Author = "Benjamin Pannell"
	app.Email = "admin@sierrasoftworks.com"
	app.Copyright = "Sierra Softworks © 2018"
	app.Version = version

	app.Commands = cli.Commands{}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level",
			Usage: "DEBUG|INFO|WARN|ERROR",
			Value: "INFO",
		},
		cli.StringFlag{
			Name: "config,c",
			Usage: "FILE",
			Value: "config.yaml",
		}
	}

	app.Before = func(c *cli.Context) error {
		log.WithFields(log.Fields{
			"log-level": c.GlobalString("log-level"),
		}).Info("Starting")

		logLevel := c.GlobalString("log-level")
		switch strings.ToUpper(logLevel) {
		case "DEBUG":
			log.SetLevel(log.DebugLevel)
			raven = raven.With(sentry.Level(sentry.Debug))
		case "WARN":
			log.SetLevel(log.WarnLevel)
			raven = raven.With(sentry.Level(sentry.Warning))
		case "ERROR":
			log.SetLevel(log.ErrorLevel)
			raven = raven.With(sentry.Level(sentry.Error))
		default:
			log.SetLevel(log.InfoLevel)
			raven = raven.With(sentry.Level(sentry.Info))
		}

		return nil
	}

	app.Action = func(c *cli.Context) error {
		loadedConfig, err := config.ReadConfig(c.GlobalString("config"))
		if err != nil {
			return err
		}

		a, err := heimdall.NewAgent(loadedConfig)
		if err != nil {
			return err
		}

		return a.Run()
	}

	app.ExitErrHandler = func(c *cli.Context, err error) {
		raven.Capture(
			sentry.ExceptionForError(err),
		).Wait()

		cli.HandleExitCoder(err)
	}

	app.RunAndExitOnError()
}
