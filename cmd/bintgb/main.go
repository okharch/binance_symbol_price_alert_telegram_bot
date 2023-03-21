package main

import (
	"context"
	"github.com/binance-exchange/go-binance"
	gklog "github.com/go-kit/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
)

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

func main() {
	// Replace with your Binance API key and secret

	// token for t.me/binance_symbols_alerts_bot
	botToken := os.Getenv("BIN_TGB_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	//if err != nil {
	//	log.Panic(err)
	//}
	//_ = bot
	//}
	//
	//
	//	bot.Debug = true
	//
	//	log.Printf("Authorized on account %s", bot.Self.UserName)
	//
	//	u := tgbotapi.NewUpdate(0)
	//	u.Timeout = 60
	//
	//	updates, err := bot.GetUpdatesChan(u)
	//
	//	for update := range updates {
	//		if update.Message == nil { // ignore non-messages
	//			continue
	//		}
	//
	//		switch update.Message.Command() {
	//		case "track":
	//			args := strings.Fields(update.Message.Text)[1:]
	//			if len(args) != 3 {
	//				reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /track SYMBOL THRESHOLD_LOW THRESHOLD_HIGH")
	//				bot.Send(reply)
	//				continue
	//			}
	//
	//			symbol := strings.ToUpper(args[0])
	//			thresholdLow, err := strconv.ParseFloat(args[1], 64)
	//			if err != nil {
	//				reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid threshold_low")
	//				bot.Send(reply)
	//				continue
	//			}
	//			thresholdHigh, err := strconv.ParseFloat(args[2], 64)
	//			if err != nil {
	//				reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid threshold_high")
	//				bot.Send(reply)
	//				continue
	//			}
	//
	//			go func() {
	//				for {
	//					ticker, err := binanceClient.NewListPricesService().Symbol(symbol).Do()
	//					if err != nil {
	//						log.Println(err)
	//						continue
	//					}
	//
	//					price, err := strconv.ParseFloat(ticker[0].Price, 64)
	//					if err != nil {
	//						log.Println(err)
	//						continue
	//					}
	//
	//					if price <= thresholdLow {
	//						reply := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s price is below %f", symbol, thresholdLow))
	//						bot.Send(reply)
	//					} else if price >= thresholdHigh {
	//						reply := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s price is above %f", symbol, thresholdHigh))
	//						bot.Send(reply)
	//					}
	//
	//					// wait for 1 minute before checking again
	//					time.Sleep(time.Minute)
	//				}
	//			}()
	//		default:
	//			reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
	//			bot.Send(reply)
	//		}
}
