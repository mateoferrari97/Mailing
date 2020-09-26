package main

import (
	"os"

	"github.com/mateoferrari97/Kit/web"
	"github.com/mateoferrari97/Mailing/cmd/app/internal"
	"github.com/mateoferrari97/Mailing/cmd/app/internal/smtp"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	server := web.NewServer()

	handler := internal.NewHandler(server)

	config := smtp.Config{
		Host:     "smtp.gmail.com",
		Port:     "587",
		Username: os.Getenv("EMAIL"),
		Password: os.Getenv("PASSWORD"),
	}

	client, err := smtp.NewClient(config)
	if err != nil {
		return err
	}

	service := internal.NewService(client)

	handler.Ping()
	handler.RouteSendEmail(service.SendEmail)

	return server.Run(":8082")
}
