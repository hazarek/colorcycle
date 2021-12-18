package main

import (
	"hazarek/colorcycle/anim"
	"hazarek/colorcycle/utils"
)

func main() {
	img := utils.Open("assets/input/eye.png")
	heightmap := utils.Open("assets/input/heightmap.png")
	opt := anim.Options{
		GradientRepeat: 5,
		ShiftPalette:   3,
		RgbMode:        5,
	}
	an := anim.MakeAnim(img, heightmap, opt)
	an.CalculateFrames()
	an.Save("assets/output/out.png", 15)
}
