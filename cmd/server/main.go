package main

import (
	"flag"
	"os"

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

	server.Start()
}
