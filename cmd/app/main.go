package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"song_lib/internal/app"
	"song_lib/internal/config"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	path := ""
	flag.StringVar(&path, "c", "", "path to the configuration file")
	flag.Parse()

	ctx := context.Background()
	log := logrus.New()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := app.NewApp(ctx, cfg, log)

	go func() {
		app.Start()
	}()
	log.Info("server is running")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	app.Stop(ctx)
	log.Info("server shutdown")
}
