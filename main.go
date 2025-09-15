package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"goyave.dev/goyave/v5"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatal("db connect error:", err)
	}

	opts := goyave.Options{}
	server, err := goyave.New(opts)
	if err != nil {
		log.Fatal("server init error:", err)
	}

	server.RegisterRoutes(createRoutes(pool))

	server.Logger.Info("Registering hooks")
	server.RegisterSignalHook()

	server.RegisterStartupHook(func(s *goyave.Server) {
		s.Logger.Info("Server is listening", "host", s.Host())
	})

	server.RegisterShutdownHook(func(s *goyave.Server) {
		s.Logger.Info("Server is shutting down")
	})

	if err := server.Start(); err != nil {
		log.Fatal("server start error:", err)
	}
}
