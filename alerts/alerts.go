package alerts

import (
	"database/sql"
	"encoding/json"
	"github.com/binance-exchange/go-binance"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

func UpdatePrices(db *sql.DB, prices []*binance.PriceTicker) error {
	// Convert prices to JSON
	jsonData, err := json.Marshal(prices)
	if err != nil {
		return err
	}

	// Prepare stored procedure call
	stmt, err := db.Prepare("select * from update_prices($1)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Call stored procedure with JSON argument
	_, err = stmt.Exec(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Set up database connection
	db, err := sql.Open("postgres", "postgresql://username:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create Telegram bot client
	bot, err := tgbotapi.NewBotAPI("YOUR_TELEGRAM_BOT_TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	// Set up message handling
	updates := bot.ListenForWebhook("/")
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Parse alert command
		if strings.HasPrefix(update.Message.Text, "/alert ") {
			parts := strings.Split(update.Message.Text, " ")
			if len(parts) != 4 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid alert command. The format should be SYMBOL KIND PRICE.")
				bot.Send(msg)
				continue
			}
			symbol := strings.ToUpper(parts[1])
			kind, err := strconv.Atoi(parts[2])
			if err != nil || (kind != 0 && kind != 1) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid alert kind. The kind should be 0 for price goes up and 1 for price goes down.")
				bot.Send(msg)
				continue
			}
			price, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid alert price.")
				bot.Send(msg)
				continue
			}

			// Add alert to database
			userId := update.Message.From.ID
			_, err = db.Exec("INSERT INTO alerts (user_id, symbol, price, kind) VALUES ($1, $2, $3, $4)", userId, symbol, price, kind)
			if err != nil {
				log.Printf("Error adding alert to database: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error adding alert to database.")
				bot.Send(msg)
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Alert added.")
			bot.Send(msg)
		}
	}
}
