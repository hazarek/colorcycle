package utils

import (
	"bytes"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/kettek/apng"
	"github.com/mazznoer/colorgrad"
	"github.com/sizeofint/webpanimation"
)

func Open(f string) image.Image {
	img, err := imgio.Open(f)
	if err != nil {
		panic(err)
	}
	return img
}

// SaveWebp saves animated webp
func SaveWebp(filename string, frames []*image.NRGBA, delay int) {
	w, h := frames[0].Bounds().Dx(), frames[0].Bounds().Dy()
	var buf bytes.Buffer
	var err error
	webpanim := webpanimation.NewWebpAnimation(w, h, 0)
	webpanim.WebPAnimEncoderOptions.SetKmin(9)
	webpanim.WebPAnimEncoderOptions.SetKmax(17)
	defer webpanim.ReleaseMemory() // don't forget call this or you will have memory leaks
	webpConfig := webpanimation.NewWebpConfig()
	webpConfig.SetLossless(1)
	timeline := 0
	for _, frame := range frames {
		err = webpanim.AddFrame(frame, timeline, webpConfig)
		if err != nil {
			log.Fatal(err)
		}
		timeline += delay
	}
	err = webpanim.Encode(&buf) // encode animation and write result bytes in buffer
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(filename, buf.Bytes(), 0777) // write bytes on disk
}

// ApngSave save animated png
func ApngSave(f string, frames []*image.NRGBA, fps int) {
	var apngNrgba = apng.APNG{Frames: make([]apng.Frame, len(frames), len(frames))}
	for i := range frames {
		// apngNrgba.Frames[i].Image = ImageToNRGBA(frames[i])
		apngNrgba.Frames[i].Image = frames[i]
		apngNrgba.Frames[i].DelayNumerator = uint16(1)
		apngNrgba.Frames[i].DelayDenominator = uint16(fps)
	}

	file, err := os.Create(f)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	apng.Encode(file, apngNrgba)
}

//ImageToGray Converting image to grayscale
func ImageToGray(img image.Image) *image.Gray {
	grayImg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			R, G, B, _ := img.At(x, y).RGBA()
			//Luma: Y = 0.2126*R + 0.7152*G + 0.0722*B
			Y := (0.2126*float64(R) + 0.7152*float64(G) + 0.0722*float64(B)) * (255.0 / 65535)
			grayPix := color.Gray{uint8(Y)}
			grayImg.Set(x, y, grayPix)
		}
	}
	return grayImg
}

// ImageToNRGBA convert image.Image to image.NRGBA
func ImageToNRGBA(im image.Image) *image.NRGBA {
	nim := image.NewNRGBA(im.Bounds())
	for y := im.Bounds().Min.Y; y < im.Bounds().Max.Y; y++ {
		for x := im.Bounds().Min.X; x < im.Bounds().Max.X; x++ {
			nim.Set(x, y, im.At(x, y))
		}
	}
	return nim
}

// MakeRgbGradient make seamlessly RGB repeated gradient.
// 	Blend Modes;
// 	0 BlendHcl
// 	1 BlendHsv
// 	2 BlendLab
// 	3 BlendLinearRgb
// 	4 BlendLuv
// 	5 BlendRgb
// 	6 BlendOklab
func MakeRgbGradient(repeat int, blendMode int) colorgrad.Gradient {
	bMode := colorgrad.BlendMode(blendMode)
	rgbtemp := []string{"#FF0000", "#00FF00", "#0000FF"}
	rgb := []string{"#FF0000", "#00FF00", "#0000FF"}
	for i := 0; i < repeat; i++ {
		rgb = append(rgb, rgbtemp...)
	}
	rgb = append(rgb, rgb[0])
	grad, _ := colorgrad.NewGradient().
		HtmlColors(rgb...).
		Mode(bMode).
		Build()
	return grad
}

// MapRange map range to another range
func MapRange(v, v1, v2, min, max float64) float64 {
	return min + ((max-min)/(v2-v1))*(v-v1)
}

// ShiftPalette roll palette
func ShiftPalette(arr []color.Color, k int) []color.Color {
	if k < 0 || len(arr) == 0 {
		return arr
	}
	r := len(arr) - k%len(arr)
	arr = append(arr[r:], arr[:r]...)
	return arr
}
