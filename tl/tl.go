package tl

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	Layout = "2006-01-02T15:04:05-0700"
	Format = "2006-01-02"
)

type Impl struct {
	User, Token string
	Bot         *tgbotapi.BotAPI
}

func NewTL(bot *tgbotapi.BotAPI, user, token string) *Impl {
	return &Impl{
		User:  user,
		Token: token,
		Bot:   bot,
	}
}

func (t *Impl) SendMessage(message string, chatID int64, replyMessageID int) {
	msg := tgbotapi.NewMessage(chatID, message)
	if replyMessageID != 0 {
		msg.ReplyToMessageID = replyMessageID
	}
	_, err := t.Bot.Send(msg)
	if err != nil {
		log.Println(err.Error())
	}
}

func (t *Impl) LogToday(u tgbotapi.Update) {
	if !isWeekday() {
		t.SendMessage("Not a week day.", u.Message.Chat.ID, u.Message.MessageID)
		return
	}

	t.SendMessage("Logging started.", u.Message.Chat.ID, 0)

	if !isAfterSixPMInIST() {
		t.SendMessage("Hey! Logging is allowed only after 6 PM.", u.Message.Chat.ID, 0)
		return
	}

	err := t.processLogging()
	if err != nil {
		t.SendMessage(err.Error(), u.Message.Chat.ID, 0)
		return
	}

	t.SendMessage("Congrats! Logging success.", u.Message.Chat.ID, 0)
}
