package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"./config"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lazywei/go-opencv"
)

var conf  = config.GetConf()

func processImg() {

}

func downloadFile(filepath string, url string) (*os.File, error) {

	// Create the file
	out, err := os.Create(conf.TmpPath + filepath + ".jpg")
	if err != nil  {
		return out, err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return out, err
	}

	return out, nil
}

func main() {
	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = conf.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		photo:= *update.Message.Photo
		if photo != nil {
			imgId:= photo[len(photo)-1].FileID
			url, err:=bot.GetFileDirectURL(imgId)
			if err != nil {
				log.Panic(err)
				continue
			}
			_, fileErr := downloadFile(imgId, url)
			if fileErr != nil {
				log.Panic(fileErr)
				continue
			}
		}

		log.Printf("[%s] %s %s", update.Message.From.UserName, update.Message.Text, (*update.Message.Photo)[0].FileID)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}