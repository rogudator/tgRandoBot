package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("еще!"),
	),
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	token, found := os.LookupEnv("ID_TOKEN")
	if found {
		log.Println("Token is found")
	} else {
		log.Println("Token is not found")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	channelName := ""
	maxIdMessage := 0
	link := "t.me/%s/%d"

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		var msg tgbotapi.MessageConfig
		if update.Message.ForwardFromChat != nil {
			channelName = update.Message.ForwardFromChat.UserName
			maxIdMessage = update.Message.ForwardFromMessageID
			txt := fmt.Sprintf(link, channelName, rand.Intn(maxIdMessage))
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, txt)
			msg.ReplyMarkup = numericKeyboard
		} else if update.Message.Text == "еще!" {
			txt := fmt.Sprintf(link, channelName, rand.Intn(maxIdMessage))
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, txt)
		} else {
			//msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Чтобы получить случайный пост с канала, сделай форвард из него в этот бот.")
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
		updates.Clear()
	}
}
