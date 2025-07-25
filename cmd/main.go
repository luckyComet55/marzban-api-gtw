package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	pcl "github.com/luckyComet55/marzban-api-gtw/internal/panel_client"
	"github.com/sethvargo/go-envconfig"
)

type AppConfig struct {
	Username       string `env:"ADMIN_USERNAME, required"`
	Password       string `env:"ADMIN_PASSWORD, required"`
	MarzbanBaseUrl string `env:"BASE_URL, required"`
}

func main() {
	_, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found: %v", err)
	}

	var c AppConfig
	envconfig.MustProcess(context.Background(), &c)

	log.Printf("username: %s\npassword: %s\nurl: %s\n", c.Username, c.Password, c.MarzbanBaseUrl)

	cli := pcl.NewMarzbanPanelClient(c.MarzbanBaseUrl, c.Username, c.Password)

	users, err := cli.GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	for i, user := range users {
		log.Printf("user #%d %s (%s)\n", i, user.Username, user.Status)
	}
}
