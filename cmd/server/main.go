package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lepidoptera/lepidoptera/internal/mailer"
	"github.com/lepidoptera/lepidoptera/internal/server"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	resendKey := os.Getenv("RESEND_API_KEY")
	fromEmail := os.Getenv("FROM_EMAIL")
	tokenSecret := os.Getenv("TOKEN_SECRET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	m := mailer.NewResendMailer(resendKey, fromEmail)

	srv, err := server.New(dbURL, m, tokenSecret)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	// start notification worker in background
	ctx := context.Background()
	go srv.Worker.Start(ctx)

	log.Printf("lepidoptera starting on :%s", port)
	log.Fatal(srv.Start(":" + port))
}
