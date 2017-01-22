// http://docs.opencv.org/master/d0/d86/tutorial_py_image_arithmetics.html#gsc.tab=0

package main

import (
	//"fmt"
	"github.com/constabulary/gb/testdata/src/e"
	"github.com/lazywei/go-opencv/opencv"
)

type maskSize struct {
	Height int
	Width  int
}

func main() {
	//_, currentfile, _, _ := runtime.Caller(0)
	image := opencv.LoadImage("tmp/AgADAgAD76cxG4hdKEg3ZdzuUG6zRmMvSw0ABHkvmCDY_4B8QvEDAAEC.jpg")
	pepe := opencv.LoadImage("img/Sad_Pepe.png")
	cascade := opencv.LoadHaarClassifierCascade("tmp/haarcascade_frontalface_alt.xml")
	faces := cascade.DetectObjects(image)
	//res := opencv.CreateImage

	for _, value := range faces {
		// resize mask to face size
		resizedMask := opencv.Resize(pepe, value.Width(), value.Height(), opencv.CV_INTER_AREA)
		// Prepare ROI for image
		roiRect := opencv.NewRect(value.X(), value.Y(), value.Width(), value.Height())
		// some work for grey mask
		greyPepe := resizedMask.Clone()
		greyMask := resizedMask.Clone()
		maskInv := resizedMask.Clone()
		opencv.CvtColor(resizedMask, greyPepe, opencv.CV_BGR2GRAY)
		opencv.Threshold(greyPepe, greyMask, 10, 255, opencv.CV_THRESH_BINARY)
		opencv.Not(greyMask, maskInv)

		//rroi:= image.GetROI()
		//image_bg = image.Clone()
		image.SetROI(roiRect)
		//opencv.AndWithMask()
		//mask := new(maskSize{})
		opencv.Rectangle(image,
			opencv.Point{value.X() + value.Width(), value.Y()},
			opencv.Point{value.X(), value.Y() + value.Height()},
			opencv.ScalarAll(255.0), 1, 1, 0)
	}

	win := opencv.NewWindow("Face Detection")
	win.ShowImage(image)
	opencv.WaitKey(0)
}
