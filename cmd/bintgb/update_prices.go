package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/binance-exchange/go-binance"
	gklog "github.com/go-kit/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"time"
)

func updatePrices(ctx context.Context, bot *tgbotapi.BotAPI, db *sql.DB) error {
	// Check if there are any live alerts set up
	var count int
	if err := db.QueryRow("SELECT count(*) FROM alerts").Scan(&count); err != nil {
		return fmt.Errorf("error checking for live alerts: %w", err)
	}
	if count == 0 {
		fmt.Println("No live alerts set up, exiting")
		return nil
	}

	// Retrieve prices from Binance
	prices, err := GetPrices(ctx)
	if err != nil {
		return fmt.Errorf("failed to obtain prices from Binance: %w", err)
	}

	// Convert the prices map to JSON
	pricesJSON, err := json.Marshal(prices)
	if err != nil {
		return fmt.Errorf("error marshaling prices to JSON: %w", err)
	}

	// Retrieve alerts from the database using a PL/pgSQL function
	rows, err := db.Query("SELECT * FROM update_prices($1::json)", pricesJSON)
	if err != nil {
		return fmt.Errorf("error querying alerts from database: %w", err)
	}
	defer rows.Close()

	// Iterate over the alert rows and send them to the Telegram chatbot
	count = 0
	for rows.Next() {
		var userID int64
		var symbol string
		var price float64
		var kind int
		count++
		if err := rows.Scan(&userID, &symbol, &price, &kind); err != nil {
			return fmt.Errorf("error scanning alert row: %w", err)
		}

		// Compose the message to send to the chatbot
		kindStr := "down"
		if kind == 1 {
			kindStr = "up"
		}
		msgText := fmt.Sprintf("Alert: %s price went %s to %f", symbol, kindStr, price)
		msg := tgbotapi.NewMessage(userID, msgText)

		// Send the message using the Telegram bot API
		if _, err := bot.Send(msg); err != nil {
			return fmt.Errorf("error sending alert message to Telegram: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over alert rows: %w", err)
	}
	log.Printf("After requesting prices %d alerts has been sent", count)

	return nil
}

// continueUpdatePrices is a function that continually updates prices and alerts every minute or until the context expires.
func continueUpdatePrices(ctx context.Context, bot *tgbotapi.BotAPI, db *sql.DB) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := updatePrices(ctx, bot, db); err != nil {
				log.Printf("error updating prices and alerts: %s", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func initBinance(ctx context.Context) binance.Binance {
	var logger gklog.Logger
	logger = gklog.NewLogfmtLogger(gklog.NewSyncWriter(os.Stderr))
	logger = gklog.With(logger, "time", gklog.DefaultTimestampUTC, "caller", gklog.DefaultCaller)

	apiSecret := os.Getenv("BINANCE_SECRET")

	hmacSigner := &binance.HmacSigner{
		Key: []byte(apiSecret),
	}
	apiKey := os.Getenv("BINANCE_API_KEY")
	// use second return value for cancelling request when shutting down the app

	binanceService := binance.NewAPIService(
		"https://www.binance.com",
		apiKey,
		hmacSigner,
		logger,
		ctx,
	)
	return binance.NewBinance(binanceService)
}

func GetPrices(ctx context.Context) (result map[string]float64, err error) {
	binanceClient := initBinance(ctx)
	prices, err := binanceClient.TickerAllPrices()
	if err != nil {
		return
	}
	result = make(map[string]float64, len(prices))
	for _, p := range prices {
		result[p.Symbol] = p.Price
	}
	return
}
