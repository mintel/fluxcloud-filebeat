package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/mintel/fluxcloud-filebeat/pkg/config"
	"github.com/mintel/fluxcloud-filebeat/pkg/handler"
	"github.com/mintel/fluxcloud-filebeat/pkg/server"
	"log"
	"os"
	"os/signal"
)

func main() {
	cfg := config.Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalln("failed to parse config", err)
	}

	handler, err := handler.NewHandler(&cfg)
	if err != nil {
		log.Fatalln("failed to initialize Filebeat handler", err)
	}

	w := server.NewServer(cfg.Port, handler)
	if err := w.Start(); err != nil {
		log.Fatalln("failed to start webhook server", err)
	}
	log.Println("started forwarding flux event to {}", cfg.FileBeatAddress)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	if err := w.Close(); err != nil {
		log.Fatalln("error occurred while shutting down the webhook server", err)
	}
	log.Println("bye")
}
