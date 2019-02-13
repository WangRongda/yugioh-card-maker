package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
)

// Config Read from config.json
type Config struct {
	Text  map[string]texter
	Image map[string]imager
}

var colorStr = map[string]uint16{
	"White": 0xffff,
	"Black": 0,
}

type texter struct {
	X        int
	Y        int
	FontSize float64 //px
	FontPath string  //font style file
	FG       string
	Content  string
}
type imager struct {
	X    int
	Y    int
	Path string
}

func errPanic(err error, args ...interface{}) {
	if nil != err {
		fmt.Println(args...)
		panic(err)
	}
}

func main() {
	jsonBytes, err := ioutil.ReadFile("config.json")
	errPanic(err)
	config := Config{}
	errPanic(json.Unmarshal(jsonBytes, &config))

	if len(os.Args) < 2 {
		panic("Need 1 argument: card template file path")
	}
	cardTemplateFd, err := os.Open(os.Args[1])
	errPanic(err)
	defer cardTemplateFd.Close()
	cardTemplateImage, _, err := image.Decode(cardTemplateFd)
	errPanic(err)
	// cardTemplateRGBA := image.NewRGBA(image.Rectangle{image.Point{0, 0}, cardTemplateImage.Bounds().Size()})
	cardTemplateRGBA, ok := cardTemplateImage.(*image.RGBA)
	if !ok {
		panic("image to rgba fail")
	}
	for _, text := range config.Text {
		addLabel(cardTemplateRGBA, text)
	}
	for _, subImg := range config.Image {
		appendImg(cardTemplateRGBA, subImg)
	}

	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, cardTemplateRGBA)
	errPanic(err)
	err = b.Flush()
	errPanic(err)
	fmt.Println("Wrote out.png OK.")
}

func addLabel(img *image.RGBA, text texter) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(text.FontPath)
	errPanic(err, text)
	f, err := freetype.ParseFont(fontBytes)
	errPanic(err)
	c := freetype.NewContext()
	// draw.Draw(img, img.Bounds(), img, image.ZP, draw.Src)
	c.SetFont(f)
	c.SetFontSize(text.FontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	fg, ok := colorStr[text.FG]
	if !ok {
		panic("FG must White or Black")
	}
	c.SetSrc(image.NewUniform(color.Gray16{fg}))
	// c.SetHinting(font.HintingNone)
	pt := freetype.Pt(text.X, text.Y+int(c.PointToFixed(text.FontSize)>>6))
	_, err = c.DrawString(text.Content, pt)
	errPanic(err)
}

func appendImg(rgba *image.RGBA, subImg imager) {
	imageFD, err := os.Open(subImg.Path)
	errPanic(err)
	subImgImage, _, err := image.Decode(imageFD)
	errPanic(err, subImg)
	sp2 := image.Point{rgba.Bounds().Dx(), 0}
	r2 := image.Rectangle{sp2, sp2.Add(subImgImage.Bounds().Size())}
	draw.Draw(rgba, r2, subImgImage, image.Point{0, 0}, draw.Src)
}
