package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func errPanic(err error, args ...interface{}) {
	if nil != err {
		fmt.Println(args...)
		panic(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		panic("Need 1 argument: card template file path")
	}
	config := struct {
		Text  map[string]texter
		Image map[string]imager
	}{}

	jsonBytes, err := ioutil.ReadFile("config.json")
	errPanic(err)
	errPanic(json.Unmarshal(jsonBytes, &config))

	card := Carder{newTexts: config.Text, newImages: config.Image}
	card.UseTemplate(os.Args[1])
	card.Flush()
	card.ExportPNG("out.png")
}
