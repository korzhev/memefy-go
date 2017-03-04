package memefy

import (
	"path/filepath"
	"bot/config"
	"github.com/lazywei/go-opencv/opencv"
)

var conf = config.GetConf()

var pepe = opencv.LoadImage(conf.Img)
var pepeMask = opencv.LoadImage(conf.Mask)

func FaceChange(filename string) string {
	var cascade = opencv.LoadHaarClassifierCascade("haarcascade_frontalface_alt.xml")
	var img = opencv.LoadImage(filename)
	faces := cascade.DetectObjects(img)
	// fix for unexpected faces changes, may be bug
	f := make([]opencv.Rect, len(faces))
	for i, v := range faces {
		f[i] = *v
	}
	for _, face := range f {
		size := face.Width()
		resizeMask := opencv.Resize(pepe, size, size, opencv.CV_INTER_AREA)
		resizeBlackMask := opencv.Resize(pepeMask, size, size, opencv.CV_INTER_AREA)
		bg := opencv.CreateImage(img.Width(), img.Height(), img.Depth(), img.Channels())

		bg.SetROI(face)
		img.SetROI(face)

		opencv.Subtract(img, resizeBlackMask, bg)
		opencv.Add(bg, resizeMask, img)

		img.ResetROI()
		bg.ResetROI()
		bg.Release()
	}
	dir, file := filepath.Split(filename)
	newFilePath := dir + "memed_" + file
	opencv.SaveImage(newFilePath, img, 100)
	img.Release()

	return newFilePath
}
