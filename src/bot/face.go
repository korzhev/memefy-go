package main

import (
	"github.com/lazywei/go-opencv/opencv"
)

var image = opencv.LoadImage("tmp/hd_cb7b4eb443.jpg")
var pepe = opencv.LoadImage("img/pepe.png")
var cascade = opencv.LoadHaarClassifierCascade("haarcascade_frontalface_alt.xml")

func main() {
	defer image.Release()
	defer pepe.Release()
	faces := cascade.DetectObjects(image)
	// fix for unexpected faces changes, may be bug
	f := make([]opencv.Rect, len(faces))
	for i, v := range faces {
		f[i] = *v
	}
	for _, face := range f {
		size := face.Width()
		resizedMask := opencv.Resize(pepe, size, size, opencv.CV_INTER_AREA)
		image.SetROI(face)

		greyPepe := opencv.CreateImage(size, size, opencv.IPL_DEPTH_8U, 1)
		defer greyPepe.Release()

		greyMask := opencv.CreateImage(size, size, opencv.IPL_DEPTH_8U, 1)
		defer greyMask.Release()

		maskInv := opencv.CreateImage(size, size, opencv.IPL_DEPTH_8U, 1)
		defer maskInv.Release()

		opencv.CvtColor(resizedMask, greyPepe, opencv.CV_BGR2GRAY)
		opencv.Threshold(greyPepe, greyMask, 10, 255, opencv.CV_THRESH_BINARY)
		opencv.Not(greyMask, maskInv)

		bg := opencv.CreateImage(image.Width(), image.Height(), image.Depth(), image.Channels())
		defer bg.Release()
		bg.SetROI(face)
		opencv.AndWithMask(image, image, bg, maskInv)

		fg := opencv.CreateImage(resizedMask.Width(), resizedMask.Height(), resizedMask.Depth(), resizedMask.Channels())
		defer fg.Release()
		opencv.AndWithMask(resizedMask, resizedMask, fg, greyPepe)

		opencv.Add(bg, fg, image)
		image.ResetROI()
		bg.ResetROI()
	}
	opencv.SaveImage("tmp/memed.jpg", image, 90)
	//win := opencv.NewWindow("Face Detection")
	//win.ShowImage(image)
	//opencv.WaitKey(0)
}
