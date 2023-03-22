package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

// alertsHandler function takes in three parameters: bot *tgbotapi.BotAPI, message *tgbotapi.Message, and db *sql.DB. The db parameter is passed in from the handleCommands function. The function retrieves all alerts for the current user from the database and sends them to the Telegram chatbot. If there are no live alerts set up for the user, the function sends a message to the chatbot saying so.
//
// If an error occurs while querying the database or sending a message to the chatbot, the function returns an error, which is handled by the handleCommands function. If the function completes successfully, it returns nil.
//
// If the active_since column of an alert row is NULL, the function appends (inactive) to the end of the message. Otherwise, it appends (active since <timestamp>).
func alertsHandler(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) error {
	// Query the database to retrieve all alerts for the current user
	query := `
		SELECT a.symbol, CASE WHEN a.kind = 0 THEN 'down' ELSE 'up' END, a.price, a.active_since, p.price
		FROM alerts a
		LEFT JOIN symbol_prices p ON a.symbol = p.symbol
		WHERE a.user_id = $1
		ORDER BY a.created_at ASC
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

	// Compose the message to send to the Telegram chatbot
	var alerts []string
	for rows.Next() {
		var symbol string
		var kind string
		var price float64
		var activeSince sql.NullTime
		var currentPrice sql.NullFloat64

		if err := rows.Scan(&symbol, &kind, &price, &activeSince, &currentPrice); err != nil {
			return fmt.Errorf("error scanning alert row: %w", err)
		}

		var status string
		if activeSince.Valid {
			status = fmt.Sprintf("(active since %s)", timeAgo(activeSince.Time))
		} else {
			status = "(inactive)"
		}

		if currentPrice.Valid {
			alerts = append(alerts, fmt.Sprintf("%s %s %f %s (current price %f)", symbol, kind, price, status, currentPrice.Float64))
		} else {
			alerts = append(alerts, fmt.Sprintf("%s %s %f %s", symbol, kind, price, status))
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over alert rows: %w", err)
	}

	msgText := strings.Join(alerts, "\n")
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)

	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message to Telegram: %w", err)
	}

	return nil
}
