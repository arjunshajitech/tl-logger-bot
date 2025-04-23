package main

import (
	"log"
	"os"
	"techgentsia-bot/tl"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	BotToken := os.Getenv("BOT_TOKEN")
	if BotToken == "" {
		panic("BOT_TOKEN not configured.")
	}
	TlUserName := os.Getenv("TL_USERNAME")
	if TlUserName == "" {
		panic("TL_USERNAME not configured.")
	}
	TlToken := os.Getenv("TL_API_TOKEN")
	if TlToken == "" {
		panic("TL_API_TOKEN not configured.")
	}

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	tl := tl.NewTL(TlUserName, TlToken)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.Text == "Log" {
				if !isWeekday() {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Not a week day.")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Logging started.")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)

				isLoggedToday, err := tl.LoggedToday()
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					bot.Send(msg)
					continue
				}

				if isLoggedToday {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Already logged today.")
					bot.Send(msg)
					continue
				} else {
					ist, err := time.LoadLocation("Asia/Kolkata")
					if err != nil {
						panic(err)
					}

					now := time.Now().In(ist)
					sixPM := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())

					if !now.After(sixPM) {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hey! Logging is allowed only after 6 PM.")
						bot.Send(msg)
						continue
					}

					err = tl.LogToday()
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
						bot.Send(msg)
						continue
					}

					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Logging success.")
					bot.Send(msg)
				}
			}
		}
	}
}

func isWeekday() bool {
	weekday := time.Now().Weekday()
	return weekday != time.Saturday && weekday != time.Sunday
}
