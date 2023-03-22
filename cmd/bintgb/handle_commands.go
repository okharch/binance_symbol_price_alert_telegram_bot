package main

import (
	"context"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func unknownCommandHandler(bot *tgbotapi.BotAPI, m *tgbotapi.Message, db *sql.DB) error {
	firstName := m.From.FirstName
	lastName := m.From.LastName

	msg := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("Hello, %s %s! This is an unknown command. Please use /help to list the available commands.", firstName, lastName))
	_, err := bot.Send(msg)
	return err
}

type CmdHandler func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) error

// handleCommands is a command handler that handles all supported commands.
func handleCommands(ctx context.Context, bot *tgbotapi.BotAPI, db *sql.DB) {
	cmdHandlers := map[string]CmdHandler{
		"help":   helpHandler,
		"alert":  handleAlert,
		"alerts": alertsHandler,
	}

	// Set up an update configuration to receive updates from the bot
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Use long polling to receive updates
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			cmd := update.Message.Command()
			handler, ok := cmdHandlers[cmd]
			if !ok {
				handler = unknownCommandHandler
			}

			if err := handler(bot, update.Message, db); err != nil {
				log.Printf("Error handling command %s cmd: %s", cmd, err)
			}
		case <-ctx.Done():
			return
		}
	}
}
