package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

// alertsHandler function takes in three parameters: bot *tgbotapi.BotAPI, message *tgbotapi.Message, and db *sql.DB. The db parameter is passed in from the handleCommands function. The function retrieves all alerts for the current user from the database and sends them to the Telegram chatbot. If there are no live alerts set up for the user, the function sends a message to the chatbot saying so.
//
// If an error occurs while querying the database or sending a message to the chatbot, the function returns an error, which is handled by the handleCommands function. If the function completes successfully, it returns nil.
//
// If the active_since column of an alert row is NULL, the function appends (inactive) to the end of the message. Otherwise, it appends (active since <timestamp>).
func alertsHandler(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) error {
	// Query the database to retrieve all alerts for the current user
	query := `
		SELECT symbol, CASE WHEN kind = 0 THEN 'down' ELSE 'up' END, price, active_since
		FROM alerts
		WHERE user_id = $1
		ORDER BY created_at ASC
	`
	rows, err := db.Query(query, message.From.ID)
	if err != nil {
		return fmt.Errorf("error querying alerts from database: %w", err)
	}
	defer rows.Close()

	// If there are no rows returned, there are no live alerts set up for the user
	if !rows.Next() {
		msg := tgbotapi.NewMessage(message.Chat.ID, "No live alerts set up.")
		if _, err := bot.Send(msg); err != nil {
			return fmt.Errorf("error sending message to Telegram: %w", err)
		}
		return nil
	}

	// Iterate over the alert rows and send them to the Telegram chatbot
	for rows.Next() {
		var symbol string
		var kind string
		var price float64
		var activeSince sql.NullTime

		if err := rows.Scan(&symbol, &kind, &price, &activeSince); err != nil {
			return fmt.Errorf("error scanning alert row: %w", err)
		}

		var status string
		if activeSince.Valid {
			status = fmt.Sprintf("(active since %s)", activeSince.Time.Format(time.RFC1123))
		} else {
			status = "(inactive)"
		}

		msgText := fmt.Sprintf("%s %s %f %s", symbol, kind, price, status)
		msg := tgbotapi.NewMessage(message.Chat.ID, msgText)

		if _, err := bot.Send(msg); err != nil {
			return fmt.Errorf("error sending message to Telegram: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over alert rows: %w", err)
	}

	return nil
}
