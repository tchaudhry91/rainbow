package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/tchaudhry91/rainbow/service"

	"github.com/peterbourgon/ff/v3"
)

func main() {
	fs := flag.NewFlagSet("rainbow", flag.ExitOnError)
	var (
		listenAddr = fs.String("listen-addr", "localhost:8080", "listen address")
		dbAddr     = fs.String("db-addr", "", "Database address")
		dbPassword = fs.String("db-password", "", "Database password")
		dbNumber   = fs.Int("db-name", 1, "Database number")
	)

	ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("RAINBOW"),
	)

	store := service.NewRedisStore(*dbAddr, *dbPassword, *dbNumber)
	svc := service.NewSHA256RainbowService(store)
	server := service.NewRainbowHTTP(*listenAddr, svc)

	shutdown := make(chan error, 1)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	go func() {
		logger.Log("msg", "Starting server..", "listenAddr", *listenAddr)
		err = server.Start()
		shutdown <- err
	}()

	select {
	case signalKill := <-interrupt:
		logger.Log("msg", fmt.Sprintf("Stopping Server: %s", signalKill.String()))
	case err := <-shutdown:
		logger.Log("error", err)
	}

	err = server.Shutdown(context.TODO())
	if err != nil {
		logger.Log("error", err)
	}
}
