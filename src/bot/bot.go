package main

import (
	"io"
	"log"
	"net/http"
	"os"
	//"path"
	"bot/config"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	//"github.com/lazywei/go-opencv/opencv"
)

var conf  = config.GetConf()

//func processImg() {
//
//}

//func findFace(filepath string) {
	//image := opencv.LoadImage(filepath)
	//
	//cascade := opencv.LoadHaarClassifierCascade(path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml"))
	//faces := cascade.DetectObjects(image)
	//
	//for _, value := range faces {
	//	opencv.Rectangle(image,
	//		opencv.Point{value.X() + value.Width(), value.Y()},
	//		opencv.Point{value.X(), value.Y() + value.Height()},
	//		opencv.ScalarAll(255.0), 1, 1, 0)
	//}
	//
	//win := opencv.NewWindow("Face Detection")
	//win.ShowImage(image)
	//opencv.WaitKey(0)
//}

func downloadFile(filepath string, url string) (string, error) {
	filename := conf.TmpPath + filepath + ".jpg"

	// Create the file
	out, err := os.Create(filename)
	if err != nil  {
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
	if err != nil  {
		return filename, err
	}

	return filename, nil
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