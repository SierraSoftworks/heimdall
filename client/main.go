package main

import (
	"os"
	"strings"
	"syscall"

	"os/signal"

	"github.com/SierraSoftworks/heimdall"
	log "github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"github.com/urfave/cli"
)

func main() {
	if envDSN := os.Getenv("SENTRY_DSN"); envDSN != "" {
		raven.SetDSN(envDSN)
	}

	raven.SetRelease(heimdall.Version)
	raven.SetEnvironment("client")

	app := cli.NewApp()
	app.Name = "Heimdall Client"
	app.Usage = "Run a Heimdall check client"

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

		return nil
	}

	app.Action = func(c *cli.Context) error {
		client, err := NewClient(c.GlobalString("config"))
		if err != nil {
			log.WithError(err).Error()
			return err
		}

		// Wait for a kill/interrupt signal before exiting
		quitSig := make(chan os.Signal, 1)
		signal.Notify(quitSig, os.Interrupt, os.Kill, syscall.SIGQUIT)

		<-quitSig
		signal.Reset(os.Interrupt, os.Kill, syscall.SIGQUIT)

		client.Shutdown()

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		os.Exit(2)
	}

	os.Exit(0)
}
