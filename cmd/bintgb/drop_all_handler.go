package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// retrieves the user ID from the message parameter, then deletes all alerts for the current user using a SQL DELETE statement. Finally, it sends a confirmation message to the chat using the bot.Send method
func dropAllHandler(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) error {
	userID := message.From.ID

	// Delete all alerts for the current user
	_, err := db.Exec("DELETE FROM alerts WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("error deleting alerts: %w", err)
	}

	// Send a message to confirm the deletion
	msg := tgbotapi.NewMessage(message.Chat.ID, "All alerts deleted")
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}
