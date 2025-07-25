package main

import (
	"flag"
	"fmt"
	"log"

	pcl "github.com/luckyComet55/marzban-api-gtw/internal/panel_client"
)

func main() {
	username := flag.String("username", "", "Marzban admin username")
	password := flag.String("password", "", "Marzban admin password")
	marzbanBaseUrl := flag.String("url", "", "Marzban base url")

	flag.Parse()

	fmt.Printf("username: %s\npassword: %s\nurl: %s\n", *username, *password, *marzbanBaseUrl)

	cli := pcl.NewMarzbanPanelClient(*marzbanBaseUrl, *username, *password)

	users, err := cli.GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	for i, user := range users {
		fmt.Printf("user #%d %s (%s)\n", i, user.Username, user.Status)
	}
}
