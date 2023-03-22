package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// Create a new bot instance using your bot token
	telegramBotToken := os.Getenv("BIN_TGB_TOKEN")
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Fatal(err)
	}

	// Set up an update configuration to receive updates from the bot
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Use long polling to receive updates
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Process updates received from the bot
	log.Printf("listen to messages for bot %s", telegramBotToken)
	for update := range updates {
		// Check if the update contains a message
		if update.Message == nil {
			continue
		}

		// Check if the message contains the /hello command
		if update.Message.IsCommand() && update.Message.Command() == "hello" {
			// Send a message to the user with the "Hello, username!" response
			u := update.Message.Chat
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hello, %s %s (%s)!", u.FirstName, u.LastName, u.UserName))
			log.Printf("arrived message from %d", update.Message.Chat.ID)
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
