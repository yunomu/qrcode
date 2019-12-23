package qrcode

import (
	"image"
	"math"

	"github.com/disintegration/imaging"
	qrcode "github.com/skip2/go-qrcode"
)

type Generator struct {
	size           int
	logoOccupancy  float64
	logoMargin     int
	recoveryLevel  qrcode.RecoveryLevel
	qrcodeMargin   int
	resampleFilter imaging.ResampleFilter
}

const (
	qrcodeMarginRate = 0.13
)

func NewGenerator(ops ...Option) *Generator {
	g := &Generator{
		logoOccupancy:  0.05,
		logoMargin:     2,
		recoveryLevel:  qrcode.High,
		resampleFilter: imaging.Linear,
	}
	Size(300)(g)

	for _, op := range ops {
		op(g)
	}

	return g
}

type Option func(g *Generator)

func Size(size int) Option {
	return func(g *Generator) {
		g.size = size
		g.qrcodeMargin = int(float64(size) * qrcodeMarginRate)
	}
}

func LogoOccupancy(o float64) Option {
	return func(g *Generator) {
		g.logoOccupancy = o
	}
}

func LogoMargin(m int) Option {
	return func(g *Generator) {
		g.logoMargin = m
	}
}

// RecoveryLevel set recovery level for QR code. Default: High
// ref. github.com/skip2/go-qrcode
func RecoveryLevel(recoveryLevel qrcode.RecoveryLevel) Option {
	return func(g *Generator) {
		g.recoveryLevel = recoveryLevel
	}
}

func ResampleFilter(resampleFilter imaging.ResampleFilter) Option {
	return func(g *Generator) {
		g.resampleFilter = resampleFilter
	}
}

func calcLogoSize(w, h, l int, rate float64) (int, int) {
	area := math.Pow(float64(l), 2) * rate
	r := float64(w) / float64(h)

	reth := math.Sqrt(area / r)
	retw := reth * r

	return int(retw), int(reth)
}

// Generate generate QR Code with logo image.
func (g *Generator) Generate(content string, logo image.Image) (image.Image, error) {
	gen, err := qrcode.New(content, g.recoveryLevel)
	if err != nil {
		return nil, err
	}

	img := gen.Image(g.size)
	if logo != nil {
		p := logo.Bounds().Max
		lw, lh := calcLogoSize(p.X, p.Y, g.size-g.qrcodeMargin, g.logoOccupancy)

		resizedLogo := imaging.OverlayCenter(
			imaging.Overlay(image.NewNRGBA(image.Rect(0, 0, lw, lh)), image.White, image.ZP, 1.0),
			imaging.Resize(logo, lw-g.logoMargin, lh-g.logoMargin, g.resampleFilter),
			1.0,
		)

		img = imaging.OverlayCenter(img, resizedLogo, 1.0)
	}

	return img, nil
}
