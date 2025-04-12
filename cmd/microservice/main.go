package main

import (
	"context"
	api "microservice_t/internal/API"
	"microservice_t/internal/config"
	logging "microservice_t/internal/logger"
	"microservice_t/internal/storage"
	"microservice_t/internal/storage/hashmap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//ctx := context.Background()

	var storage storage.Storage

	var errShutdown = make(chan error, 1)
	// get config # viper

	//cfg := config.MustReadConfig()

	cfg := config.DebugConfig()

	log := logging.New(cfg.Log)

	log.Info("logger is runned")

	// storage, err := postorage.New(ctx, log, config)
	// if err != nil {
	// 	log.Error("database not initizalited", "err", err.Error())

	// 	errShutdown <- err
	// }

	storage = hashmap.New()

	log.Info("Storage is runned")

	go storage.Ping(context.TODO(), errShutdown)

	// start API - http server

	api := api.NewAPI(log, storage)

	srv := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      api.Router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimout,
	}

	go func(errShutdown chan error) {
		err := srv.ListenAndServe()

		if err != nil {
			errShutdown <- err
		}
	}(errShutdown)

	log.Info("http server is runned")

	//	graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

Loop:
	for {
		select {
		case sig := <-shutdown:
			log.Warn("Has notifited os signal. Shutting down", "sig", sig.String())
			break Loop
		case err := <-errShutdown:
			log.Error("Critical error. Shutting down", "err", err)
			break Loop
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", "err", err.Error())
	}

	storage.Close()
	log.Debug("storage is closed")

}
