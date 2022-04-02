package main

//go:generate go get github.com/rakyll/statik
//go:generate statik

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
	//_ "./statik"
	_ "github.com/Doarakko/bigburger/statik"
	"github.com/mattn/go-sixel"
	"github.com/mattn/longcat/iterm"
	"github.com/rakyll/statik/fs"
)

// Topping struct
type Topping struct {
	Name   string
	Count  int
	Option string
	Img    image.Image
}

func loadImage(fs http.FileSystem, n string) (image.Image, error) {
	f, err := fs.Open(n)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	return png.Decode(f)
}

func saveImage(filename string, img image.Image) error {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, buf.Bytes(), 0644)
}

var toppings [6]Topping
var top Topping
var topWithSesame Topping
var bottom Topping

func init() {
	fs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	top.Img, err = loadImage(fs, "/top.png")
	if err != nil {
		log.Fatal(err)
	}

	topWithSesame.Img, err = loadImage(fs, "/top-sesame.png")
	if err != nil {
		log.Fatal(err)
	}

	bottom.Img, err = loadImage(fs, "/bottom.png")
	if err != nil {
		log.Fatal(err)
	}

	toppings[0].Name = "putty"
	toppings[0].Count = 1
	toppings[0].Option = "p"

	toppings[1].Name = "cat"
	toppings[1].Count = 0
	toppings[1].Option = "C"

	toppings[2].Name = "cheese"
	toppings[2].Count = 0
	toppings[2].Option = "c"

	toppings[3].Name = "bun"
	toppings[3].Count = 0
	toppings[3].Option = "b"

	toppings[4].Name = "tomato"
	toppings[4].Count = 0
	toppings[4].Option = "t"

	toppings[5].Name = "lettuce"
	toppings[5].Count = 0
	toppings[5].Option = "l"

	for i := 0; i < len(toppings); i++ {
		toppings[i].Img, err = loadImage(fs, fmt.Sprintf("/%s.png", toppings[i].Name))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	var nBurger int
	var rInterval float64
	var fileName string
	var sesame bool

	flag.IntVar(&nBurger, "n", 1, "number of big burger")
	flag.Float64Var(&rInterval, "i", 1.0, "rate of intervals")
	flag.StringVar(&fileName, "o", "", "output image file")
	flag.BoolVar(&sesame, "s", false,"buns with sesame")

	for i := 0; i < len(toppings); i++ {
		flag.IntVar(&toppings[i].Count, toppings[i].Option, toppings[i].Count, fmt.Sprintf("how many %s", toppings[i].Name))
	}
	flag.Parse()

	if sesame{
		top = topWithSesame
	}

	// Count number of unique toppings
	uniqueToppingCount := 0
	for i := 0; i < len(toppings); i++ {
		if toppings[i].Count > 0 {
			uniqueToppingCount++
		}
	}

	width := int(float64(top.Img.Bounds().Dx()*(nBurger-1))*rInterval) + top.Img.Bounds().Dx()
	height := top.Img.Bounds().Dy() + bottom.Img.Bounds().Dy()
	for i := 0; i < len(toppings); i++ {
		height += toppings[i].Img.Bounds().Dy() * toppings[i].Count
	}
	rect := image.Rect(0, 0, width, height)
	canvas := image.NewRGBA(rect)

	for col := 0; col < nBurger; col++ {
		// top
		x := int(float64(top.Img.Bounds().Dx()*col) * rInterval)
		rect = image.Rect(x, 0, x+top.Img.Bounds().Dx(), top.Img.Bounds().Dy())
		draw.Draw(canvas, rect, top.Img, image.Point{}, draw.Over)

		// toppings
		count := 0
		y := top.Img.Bounds().Dy()
		for i := 0; count != uniqueToppingCount; i++ {
			for j := 0; j < len(toppings); j++ {
				if toppings[j].Count == 0 || toppings[j].Count-i < 0 {
					continue
				}

				if toppings[j].Count-i > 0 {
					rect = image.Rect(x, y, x+top.Img.Bounds().Dx(), y+toppings[j].Img.Bounds().Dy())
					draw.Draw(canvas, rect, toppings[j].Img, image.Point{}, draw.Over)

					y += toppings[j].Img.Bounds().Dy()
				} else {
					count++
				}
			}
		}

		// bottom
		rect = image.Rect(x, y, x+top.Img.Bounds().Dx(), y+bottom.Img.Bounds().Dy())
		draw.Draw(canvas, rect, bottom.Img, image.Point{}, draw.Over)
	}

	var output image.Image = canvas
	if fileName != "" {
		err := saveImage(fileName, output)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	var buf bytes.Buffer
	var enc interface {
		Encode(image.Image) error
	}
	if checkIterm() {
		enc = iterm.NewEncoder(&buf)
	} else {
		enc = sixel.NewEncoder(&buf)
	}

	if err := enc.Encode(output); err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(buf.Bytes())
	os.Stdout.Sync()
}

// https://github.com/mattn/longcat
func checkIterm() bool {
	if strings.HasPrefix(os.Getenv("TERM_PROGRAM"), "iTerm") {
		return true
	}
	return getDA2() == "\x1b[>0;95;0c" // iTerm2 version 3
}

// https://github.com/mattn/longcat
func getDA2() string {
	s, err := term.MakeRaw(1)
	if err != nil {
		return ""
	}
	defer term.Restore(1, s)
	_, err = os.Stdout.Write([]byte("\x1b[>c")) // DA2 host request
	if err != nil {
		return ""
	}
	defer os.Stdout.SetReadDeadline(time.Time{})

	time.Sleep(10 * time.Millisecond)

	var b [100]byte
	n, err := os.Stdout.Read(b[:])
	if err != nil {
		return ""
	}
	return string(b[:n])
}
