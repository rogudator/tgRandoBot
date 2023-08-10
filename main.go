package main

import (
	"errors"
	"os"

	// "os"

	//"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type html struct {
	Body string `xml:"body"`
}

func GetLastPostId(channelName string) (int, error) {
	url := "https://t.me/s/"
	resp, err := http.Get(url + channelName)
	if err != nil {
		return 0, err
	}

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	bodyString := string(bodyByte)
	beforeNumber := "href=\"https://t.me/" + channelName + "/"
	trimBefore := strings.LastIndex(bodyString, beforeNumber)
	lastPostString := ""

	if trimBefore != -1 {
		trimBefore += len(beforeNumber)
		bodyString = bodyString[trimBefore:]
		for _, v := range bodyString {
			if v == '"' {
				lastPostID, err := strconv.Atoi(lastPostString)
				if err != nil {
					return 0, err
				}
				return lastPostID, nil
			}
			lastPostString += string(v)
		}

	}
	return 0, errors.New("Can't get id of channel's last post.")
}

func GetRandomPostLink(channelName string) (string, error) {
	lastId, err := GetLastPostId(channelName)
	if err != nil {
		return "", err
	}
	// lastId := 14611
	randoId := rand.Intn(lastId)
	link := "t.me/" + channelName + "/" + strconv.Itoa(randoId)

	return link, nil
}

func GetRandomPostID(channelName string) (int) {
	randoId := 14
	lastId, err := GetLastPostId(channelName)
	if err != nil {
		log.Println(err)
		return rand.Intn(randoId)
	}
	randoId = rand.Intn(lastId)

	return randoId
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("moreeee"),
	),
)

// init is invoked before main()
func init() {
    // loads values from .env into the system
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
}

func main() {
	token, ok := os.LookupEnv("ID_TOKEN")
	if !ok {
		log.Fatal("couldn't read id token")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	channel, ok := os.LookupEnv("CHANNEL")
	if !ok {
		log.Fatal("couldnt read channel name")
	}
	channelString, ok := os.LookupEnv("CHANNEL_ID")
	if !ok {
		log.Fatal("couldnt read channel id name")
	}
	channelID, err := strconv.ParseInt(channelString, 10, 64)
	if err != nil {
		log.Fatal("failed to convert env to valid id")
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)



	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		forwardIsNotSent := true
		for forwardIsNotSent {
			postID := GetRandomPostID(channel)
			msg2 := tgbotapi.NewForward(update.Message.Chat.ID , channelID, postID) 
			_, err := bot.Send(msg2)
			if err == nil {
				forwardIsNotSent = false
			}
		}
		var msg tgbotapi.MessageConfig
		msg.Text = "Another post?"
		msg.ChatID = update.Message.Chat.ID
		msg.ReplyMarkup = numericKeyboard
		
		bot.Send(msg)
		updates.Clear()
	}
}
