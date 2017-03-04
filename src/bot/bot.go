package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"bot/config"
	"bot/memefy"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var conf = config.GetConf()

func downloadFile(filePath string, url string) (string, error) {
	filename := conf.TmpPath + filePath + ".jpg"

	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		return filename, err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return filename, err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return filename, err
	}

	return filename, nil
}

func process(bot *tgbotapi.BotAPI, chat int64, imgId string) {
	url, err := bot.GetFileDirectURL(imgId)
	if err != nil {
		log.Fatal(err)
		return
	}

	filename, fileErr := downloadFile(imgId, url)
	if fileErr != nil {
		log.Fatal(err)
		return
	}
	changedFile := memefy.FaceChange(filename)
	errRemoveIncoming := os.Remove(filename)
	if errRemoveIncoming != nil {
		log.Fatal(errRemoveIncoming)
		return
	}
	msg := tgbotapi.NewPhotoUpload(chat, changedFile)
	msg.Caption = "@memefypepefy_bot kek"
	_, errMsg := bot.Send(msg)
	if errMsg != nil {
		log.Fatal(errMsg)
		return
	}
	errRemoveChanged := os.Remove(changedFile)
	if errRemoveChanged != nil {
		log.Fatal(errRemoveChanged)
	}
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
		photo := *update.Message.Photo
		if photo != nil {
			imgId := photo[len(photo)-1].FileID
			go process(bot, update.Message.Chat.ID, imgId)
		}
	}
}
