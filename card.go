package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
)

type texter struct {
	X        int
	Y        int
	FontSize float64 //px
	FontPath string  //font style file
	FG       string
	Content  string
}
type imager struct {
	X      int
	Y      int
	Width  int //从左截取,非缩放
	Height int //从上截取，非缩放
	Path   string
}

var colorStr = map[string]uint16{
	"White": 0xffff,
	"Black": 0,
}

type Carder struct {
	baseRGBA  *image.RGBA
	newTexts  map[string]texter
	newImages map[string]imager
}

func (card Carder) drawLabel(text texter) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(text.FontPath)
	errPanic(err, text)
	f, err := freetype.ParseFont(fontBytes)
	errPanic(err)
	c := freetype.NewContext()
	// draw.Draw(img, img.Bounds(), img, image.ZP, draw.Src)
	c.SetFont(f)
	c.SetFontSize(text.FontSize)
	c.SetClip(card.baseRGBA.Bounds())
	c.SetDst(card.baseRGBA)
	fg, ok := colorStr[text.FG]
	if !ok {
		panic(fmt.Sprintf("FG color not in %v", colorStr))
	}
	c.SetSrc(image.NewUniform(color.Gray16{fg}))
	// c.SetHinting(font.HintingNone)
	pt := freetype.Pt(text.X, text.Y+int(c.PointToFixed(text.FontSize)>>6))
	_, err = c.DrawString(text.Content, pt)
	errPanic(err)
}

func (card Carder) drawImage(subImg imager) {
	imageFD, err := os.Open(subImg.Path)
	errPanic(err)
	subImgImage, _, err := image.Decode(imageFD)
	errPanic(err, subImg)
	sp2 := image.Point{subImg.X, subImg.Y} //插入的图片的位置
	r2 := image.Rectangle{sp2, sp2.Add(image.Point{subImg.Width, subImg.Height})}
	draw.Draw(card.baseRGBA, r2, subImgImage, image.Point{0, 0}, draw.Src)
}

func (card *Carder) UseTemplate(filepath string) {
	cardTemplateFd, err := os.Open(filepath)
	errPanic(err)
	defer cardTemplateFd.Close()
	cardTemplateImage, _, err := image.Decode(cardTemplateFd)
	errPanic(err)
	// cardTemplateRGBA := image.NewRGBA(image.Rectangle{image.Point{0, 0}, cardTemplateImage.Bounds().Size()})
	baseRGBA, ok := cardTemplateImage.(*image.RGBA)
	if !ok {
		panic("image to rgba fail")
	}
	card.baseRGBA = baseRGBA
}

func (card *Carder) Flush() {
	for _, text := range card.newTexts {
		card.drawLabel(text)
	}
	for _, subImg := range card.newImages {
		card.drawImage(subImg)
	}
}

func (card *Carder) DrawText(name string, text texter) {
	card.newTexts[name] = text
	card.drawLabel(text)
}

func (card *Carder) DrawImage(name string, image imager) {
	card.newImages[name] = image
	card.drawImage(image)
}

func (card Carder) ExportPNG(filename string) {
	outFile, err := os.Create(filename)
	errPanic(err)
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, card.baseRGBA)
	errPanic(err)
	err = b.Flush()
	errPanic(err)
	fmt.Println("Wrote out.png OK.")
}
