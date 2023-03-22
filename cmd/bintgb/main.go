package main

import (
	"context"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// Retrieve the PostgreSQL database URL from the environment variable
	dbURL := os.Getenv("BINTBDB")

	// Connect to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new bot instance using your bot token
	telegramBotToken := os.Getenv("BIN_TGB_TOKEN")
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Wait for a SIGTERM signal
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM)
		<-sigChan

		// Signal the context to expire
		cancel()
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		continueUpdatePrices(ctx, bot, db)
		wg.Done()
	}()

	handleCommands(ctx, bot, db)
	wg.Wait()
}
