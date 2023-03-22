package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

// handleAlert is a command handler that adds an alert to the alerts table in the database.
// It expects a message in the format "/alert symbol down/up price".
func handleAlert(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) error {
	// Parse the message text to retrieve the alert details
	parts := strings.Split(message.Text, " ")
	if len(parts) != 4 {
		return fmt.Errorf("invalid alert command. Usage: /alert symbol down/up price")
	}

	symbol := parts[1]
	direction := parts[2]
	price, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return fmt.Errorf("invalid price: %s", err)
	}

	// Determine the alert kind (0 for down, 1 for up)
	kind := 0
	if direction == "up" {
		kind = 1
	}

	// Insert the alert into the database
	userID := message.Chat.ID
	_, err = db.Exec("SELECT add_alert($1, $2, $3, $4)", userID, symbol, kind, price)
	if err != nil {
		return fmt.Errorf("error adding alert to database: %w", err)
	}

	// Send a confirmation message to the user
	msgText := fmt.Sprintf("alert set for %s %s %f", symbol, direction, price)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("error sending confirmation message to Telegram: %s", err)
	}

	return nil
}
