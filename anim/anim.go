package anim

import (
	"fmt"
	"hazarek/colorcycle/utils"
	"image"
	"image/color"

	"github.com/nfnt/resize"
)

// Anim animation object
type Anim struct {
	InputImage     image.Image
	InputHeightmap image.Image
	options        Options
	frames         []*image.NRGBA
	gradient       []color.Color
	frameCount     int
}

// Options set anim opt
type Options struct {
	GradientRepeat int
	ShiftPalette   int
	// RgbMode modes;
	// 	Blend Modes;
	// 	0 BlendHcl
	// 	1 BlendHsv
	// 	2 BlendLab
	// 	3 BlendLinearRgb
	// 	4 BlendLuv
	// 	5 BlendRgb
	// 	6 BlendOklab
	RgbMode int
}

// SaveFrame calculates and saves frame with given index
func (a *Anim) SaveFrame(i int) *image.NRGBA {
	return utils.ImageToNRGBA(a.InputImage)
}

// Save saves animation as APNG, WEBP
func (a *Anim) Save(filename string, fps int) {
	fmt.Println("Saving...")
	utils.ApngSave(filename, a.frames, fps)
	fmt.Println("Done")
	// utils.SaveWebp("assets/output/anim.webp", frames, 1000/15)
}

// CalculateFrames calculates animation
func (a *Anim) CalculateFrames() {
	temp := utils.ImageToGray(utils.ImageToNRGBA(a.InputHeightmap))
	var ts uint32 = 0
	for i := 0; i < a.frameCount; i++ {
		frame := image.NewNRGBA(a.InputHeightmap.Bounds())
		clr := color.RGBA{}
		clr.A = 255
		for y := 0; y < a.InputHeightmap.Bounds().Dy(); y++ {
			for x := 0; x < a.InputHeightmap.Bounds().Dx(); x++ {
				// demRGB.Set(x, y, a.gradient[temp.GrayAt(x, y).Y])
				lR, lG, lB, _ := a.gradient[temp.GrayAt(x, y).Y].RGBA()
				r, g, b, _ := a.InputImage.At(x, y).RGBA()
				if lR+r > 65535-ts {
					clr.R = 255
				} else {
					clr.R = 0
				}
				if lG+g > 65535-ts {
					clr.G = 255
				} else {
					clr.G = 0
				}
				if lB+b > 65535-ts {
					clr.B = 255
				} else {
					clr.B = 0
				}
				frame.Set(x, y, clr)
			}
		}
		a.frames = append(a.frames, utils.ImageToNRGBA(resize.Resize(500, 500, frame, resize.Bicubic)))
		a.gradient = utils.ShiftPalette(a.gradient, a.options.ShiftPalette)
		fmt.Print(" ", i)
		if i == a.frameCount-1 {
			fmt.Println(" ")
		}
	}
}

// MakeAnim make color cycling animation
func MakeAnim(i, h image.Image, o Options) Anim {
	ani := Anim{
		InputImage:     i,
		InputHeightmap: h,
		options:        o,
	}
	ani.frameCount = 256 / o.GradientRepeat / ani.options.ShiftPalette
	ani.gradient = utils.MakeRgbGradient(o.GradientRepeat-1, ani.options.RgbMode).Colors(256)
	return ani
}
