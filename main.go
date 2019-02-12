package main

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

type texter struct {
	X        int
	Y        int
	fontSize int    //px
	fontPath string //font style file
	fg       *image.Uniform
}

// card template
type templater struct {
	background string //image path
	title      texter
	cardType   texter
	image      string
	describe   texter
}

func main() {
	//魔法卡
	spellCard := templater{
		background: "spell-template.png",
		title: texter{
			X:        49,
			Y:        27,
			fontSize: 60,
			fontPath: "YGODIY-Chinese.ttf",
			fg:       image.White,
		},
		cardType: texter{},
		image:    "",
		describe: texter{},
	}
	raw, err := os.Open("spell-template.png")
	if err != nil {
		panic(err)
	}
	defer raw.Close()
	img, err := png.Decode(raw)
	if err != nil {
		panic(err)
	}
	rgba, ok := img.(*image.RGBA)
	if !ok {
		panic("image to rgba fail")
	}
	if err != nil {
		panic(err)
	}
	x, y := 49, 27
	addLabel(rgba, x, y, "清晨的祝语", image.White)

	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, img)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")
}

func addLabel(img *image.RGBA, x, y int, label string, fg *image.Uniform) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile("YGODIY-Chinese.ttf")
	if err != nil {
		log.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}
	c := freetype.NewContext()
	draw.Draw(img, img.Bounds(), img, image.ZP, draw.Src)
	c.SetFont(f)
	c.SetFontSize(60)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)
	size := 60.0 // font size in pixels
	pt := freetype.Pt(x, y+int(c.PointToFixed(size)>>6))

	if _, err := c.DrawString(label, pt); err != nil {
		panic(err)
		// handle error
	}
}
