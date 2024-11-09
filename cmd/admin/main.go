package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

const appName = "admin"
const appAddr = ":8080"
const appEnv = "dev"

func main() {

	//TODO: Add support for config
	//TODO: Add support for config override

	logrus.Infof("Starting service: %s", appName)
	ctx, cancel := context.WithCancel(context.Background())

	app, err := newApp(ctx)
	if err != nil {
		logrus.Fatalf("Failed to init service '%s': %s", appName, err)
	}

	if err := app.Run(); err != nil {
		cancel()
		logrus.Errorf("Failed to run service '%s': %s", appName, err)
	}

	waitSignal()

	logrus.Warnf("Gracefully shutting down %s, press CTRL + C to force exit", appName)

	app.Shutdown(ctx)
}

func waitSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGINT)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGTERM)
	<-sigCh

	go func() {
		<-sigCh
		logrus.Print("Forced exit")
		os.Exit(1)
	}()
}
