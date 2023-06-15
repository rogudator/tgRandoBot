package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

// Setting up the button to ask for more messages
var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("еще!"),
	),
)

func main() {
	// Load the .env file, which contains the token
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	// Check if .env has ID_TOKEN field
	token, found := os.LookupEnv("ID_TOKEN")
	// Tell if it's found or not
	if found {
		log.Println("Token is found")
	} else {
		log.Fatal("Token is not found")
	}
	// Set up bot with token
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Set up listener to updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// This block is needed to generate link like t.me/{channelName}/{randomNumberBetween{0}and{MaxIdMessage}}
	link := "t.me/%s/%d"
	// This map is needed to remember which channel did user forward previously
	var userAndChannelName sync.Map
	// This map is needed to remember the id of post user forwarded previously
	var userAndmaxIdMessage sync.Map
	// In the long run, if the number of user grow too big. It will be a better
	// idea to store this information in database.

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}
		
		// msg is something we will send back to user
		var msg tgbotapi.MessageConfig
		// if user forwarded message from channel, we extract channel's @ and the id of forwarded message
		// then we generate link and put it in txt value
		if update.Message.ForwardFromChat != nil {
			userId := update.Message.Chat.ID
			channelName := update.Message.ForwardFromChat.UserName
			maxIdMessage := update.Message.ForwardFromMessageID
			userAndChannelName.Store(userId, channelName)
			userAndmaxIdMessage.Store(userId, maxIdMessage)
			txt := fmt.Sprintf(link, channelName, rand.Intn(maxIdMessage))
			msg = tgbotapi.NewMessage(userId, txt)
			msg.ReplyMarkup = numericKeyboard
		} else if update.Message.Text == "еще!" {
			// if user wants another random message, he sends "more" message and we repeat the generating of link
			txt := ""
			userId := update.Message.Chat.ID
			channelName, channelOk := userAndChannelName.Load(userId)
			maxIdMessageAny, maxIdMessageOk := userAndmaxIdMessage.Load(userId)
			if !channelOk || !maxIdMessageOk {
				txt = "Чтобы получить случайный пост с канала, сделай форвард из него в этот бот."
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, txt)
			} else {
				maxIdMessage, ok := maxIdMessageAny.(int)
				if !ok {
					log.Println("Error when type asserting the maxIdMessage")
					updates.Clear()
					continue;
				}
				txt = fmt.Sprintf(link, channelName, rand.Intn(maxIdMessage))
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, txt)
			}
		} else {
			// if he does not send forward, we kindly remind user to do so
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Чтобы получить случайный пост с канала, сделай форвард из него в этот бот.")
		}

		// after succesfully creating msg, we send it to the user
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
		// this is needed to not send two identical messages at the same time
		updates.Clear()
	}
}
