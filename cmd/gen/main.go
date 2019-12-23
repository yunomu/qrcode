package main

import (
	"flag"
	"image"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/yunomu/qrcode"
)

var (
	content  = flag.String("content", "It will b file, tomorrow.", "QRcode content")
	logoFile = flag.String("logo", "", "logo file (PNG)")
	outFile  = flag.String("o", "out.png", "outfile")
)

func init() {
	flag.Parse()
	log.SetOutput(os.Stderr)
}

func main() {
	var logoIn io.Reader
	if *logoFile == "-" {
		logoIn = os.Stdin
	} else if *logoFile != "" {
		f, err := os.Open(*logoFile)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		logoIn = f
	}
	var img image.Image
	if logoIn != nil {
		i, fmt, err := image.Decode(logoIn)
		if err != nil {
			log.Fatalln(err)
		}
		var _ = fmt

		img = i
	}

	var out io.Writer = os.Stdout
	if *outFile != "-" {
		f, err := os.Create(*outFile)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()

		out = f
	}

	gen := qrcode.NewGenerator()
	img, err := gen.Generate(*content, img)
	if err != nil {
		log.Fatalln(err)
	}

	if err := png.Encode(out, img); err != nil {
		log.Fatalln(err)
	}
}
