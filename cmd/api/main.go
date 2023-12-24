package main

import (
	"flag"
	svix "github.com/svix/svix-webhooks/go"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"github.com/JaegyuDev/Hydra/internal/database"
	"github.com/JaegyuDev/Hydra/internal/env"

	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL  string
	httpPort int
	db       struct {
		dsn string
	}
	clerk struct {
		clientToken string
	}
}

type application struct {
	config config
	db     *database.DB
	logger *slog.Logger
	wg     sync.WaitGroup
	clerk  struct {
		svixWebHook *svix.Webhook
	}
}

func run(logger *slog.Logger) error {
	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:3000")
	cfg.httpPort = env.GetInt("HTTP_PORT", 3000)
	cfg.db.dsn = env.GetString("DB_DSN", "user:pass@localhost:5432/db")
	cfg.clerk.clientToken = env.GetString("CLERK_TOKEN", "")
	flag.Parse()

	db, err := database.New(cfg.db.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// setup clerk conn
	wh, err := svix.NewWebhook(env.GetString("CLERK_WEBHOOK_SECRET", ""))

	app := &application{
		config: cfg,
		db:     db,
		logger: logger,
		clerk: struct {
			svixWebHook *svix.Webhook
		}{
			svixWebHook: wh,
		},
	}

	return app.serveHTTP()
}
