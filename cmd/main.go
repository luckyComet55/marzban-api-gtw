package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	app "github.com/luckyComet55/marzban-api-gtw/internal/app"
	"github.com/sethvargo/go-envconfig"
)

type AppConfig struct {
	Username       string `env:"ADMIN_USERNAME, required"`
	Password       string `env:"ADMIN_PASSWORD, required"`
	MarzbanBaseUrl string `env:"BASE_URL, required"`
	Port           uint64 `env:"PORT, default=8343"`
	Env            string `env:"ENV, required"`
}

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	_, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		slog.Warn(fmt.Sprintf("Warning: No .env file found: %v", err))
	}

	var c AppConfig
	envconfig.MustProcess(context.Background(), &c)

	cliLogger := configureLogger(c.Env, "marzban_http_cli")
	serverLogger := configureLogger(c.Env, "marzban_api_gtw_server")

	application := app.NewApp(
		serverLogger,
		cliLogger,
		c.Username,
		c.Password,
		c.MarzbanBaseUrl,
		c.Port,
	)

	go func() {
		application.MustStart()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
	serverLogger.Info("Gracefully stopped")
}

func configureLogger(env, componentName string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envDev:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).With("component", componentName)
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})).With("component", componentName)
	default:
		log.Fatalf("incorrect ENV type: %s\n", env)
	}

	return logger
}
