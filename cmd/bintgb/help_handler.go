package main

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func helpHandler(bot *tgbotapi.BotAPI, message *tgbotapi.Message, _ *sql.DB) error {
	helpText := `
		Here are the available commands:
		
		/alert <symbol> <up|down> <price> - Add an alert for a stock symbol when it reaches a certain price
		/alerts - List all active alerts for this chat
	`
	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	_, err := bot.Send(msg)
	return err
}
