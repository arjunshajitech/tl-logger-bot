package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"os/signal"
	"syscall"
	"techgentsia-bot/tl"
)

const (
	BotToken   string = "BOT_TOKEN"
	TLUserName string = "TL_USERNAME"
	TLApiToken string = "TL_API_TOKEN"
)

func main() {
	log.Println("Server started...")
	checkContainsValidEnvironment()
	go initServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigChan:
			log.Printf("Received signal: %v. Exiting...", sig)
			return
		}
	}
}

func initServer() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv(BotToken))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false

	log.Printf("Authorized telegram account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	logger := tl.NewTL(bot, os.Getenv(TLUserName), os.Getenv(TLApiToken))
	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "Log":
				logger.LogToday(update)
			case "Hello":
				logger.SendMessage("Hello there! Yes, I’m alive and listening. How can I help you today?", update.Message.Chat.ID, update.Message.MessageID)
			case "/help":
				msg := `Available Commands:
					1. Log – Log your activity for today.
					2. Logc – Log today’s activity with a commit message from the project you configured.
					3. Hello – Check if I’m alive!`
				logger.SendMessage(msg, update.Message.Chat.ID, update.Message.MessageID)
			default:
				logger.SendMessage("Oops! I didn’t recognize that command. Try /help to see what I can do!", update.Message.Chat.ID, update.Message.MessageID)
			}
		}
	}
}

func checkContainsValidEnvironment() {
	requiredVars := []string{BotToken, TLUserName, TLApiToken}
	for _, env := range requiredVars {
		if os.Getenv(env) == "" {
			panic(fmt.Sprintf("%s not configured.", env))
		}
	}
}
