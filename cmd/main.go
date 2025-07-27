package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	cfg "github.com/luckyComet55/marzban-api-gtw/internal/config"
	pcl "github.com/luckyComet55/marzban-api-gtw/internal/panel_client"
	"github.com/sethvargo/go-envconfig"
)

type AppConfig struct {
	Username       string `env:"ADMIN_USERNAME, required"`
	Password       string `env:"ADMIN_PASSWORD, required"`
	MarzbanBaseUrl string `env:"BASE_URL, required"`
	Port           int    `env:"PORT, default=8343"`
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

	logger := configureLogger(c.Env)

	log.Printf("username: %s | password: %s | url: %s", c.Username, c.Password, c.MarzbanBaseUrl)

	cli := pcl.NewMarzbanPanelClient(cfg.MarzbanApiGtwConfig{
		Username:       c.Username,
		Password:       c.Password,
		MarzbanBaseUrl: c.MarzbanBaseUrl,
		Port:           c.Port,
	}, logger)

	users, err := cli.GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	for i, user := range users {
		logger.Debug(fmt.Sprintf("user #%d %s (%s)\n", i, user.Username, user.Status))
	}
}

func configureLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envDev:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log.Fatalf("incorrect ENV type: %s\n", env)
	}

	return logger
}
