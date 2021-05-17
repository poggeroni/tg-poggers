package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func getTriggers() []string {
	triggerFile, _ := os.Open("triggers.txt")
	defer triggerFile.Close()
	sc := bufio.NewScanner(triggerFile)
	var triggers []string
	for sc.Scan() {
		triggers = append(triggers, strings.ToLower(sc.Text()))
	}
	return triggers
}

func doesTrigger(message string, triggers []string) bool {
	message = strings.ToLower(message)
	for _, trigger := range triggers {
		if strings.Contains(message, trigger) {
			return true
		}
	}
	return false
}

func getApiKey() string {
	keyFile, err := os.Open("api_key.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer keyFile.Close()

	sc := bufio.NewScanner(keyFile)
	if !sc.Scan() {
		log.Panic("ERROR: API Key file is empty.")
		os.Exit(1)
	}
	return sc.Text()
}

func main() {
	bot, err := tgbotapi.NewBotAPI(getApiKey())
	if err != nil {
		log.Panic(err)
	}
	// bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	triggers := getTriggers()

	for update := range updates {
		log.Printf("%v\n", update)
		if update.Message == nil {
			continue
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if doesTrigger(update.Message.Text, triggers) {
			msg := tgbotapi.NewStickerShare(update.Message.Chat.ID, "CAACAgQAAxkBAAMrYKJ2kZj8KG_R2lR7W23tbOfPrMYAArQJAAJZkxhRRZs6Xgc7wgABHwQ")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
