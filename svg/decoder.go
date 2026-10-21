package svg

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"math"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// Decode reads an SVG stream and rasterizes it into an *image.RGBA.
func Decode(r io.Reader) (image.Image, error) {
	icon, err := oksvg.ReadIconStream(r)
	if err != nil {
		return nil, err
	}
	if len(icon.SVGPaths) == 0 && icon.ViewBox.W == 0 && icon.ViewBox.H == 0 {
		return nil, fmt.Errorf("invalid or empty SVG: no svg element or paths found")
	}

	w := icon.ViewBox.W
	h := icon.ViewBox.H

	if w <= 0 || h <= 0 {
		w, h = 512, 512
	} else if w < 512 && h < 512 {
		scale := math.Max(512.0/w, 512.0/h)
		w = math.Round(w * scale)
		h = math.Round(h * scale)
	} else if w > 2048 || h > 2048 {
		scale := math.Min(2048.0/w, 2048.0/h)
		w = math.Round(w * scale)
		h = math.Round(h * scale)
	}

	icon.SetTarget(0, 0, w, h)
	rgba := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	scanner := rasterx.NewScannerGV(int(w), int(h), rgba, rgba.Bounds())
	dasher := rasterx.NewDasher(int(w), int(h), scanner)
	icon.Draw(dasher, 1)

	return rgba, nil
}

// DecodeConfig returns the dimensions and color model of the SVG without full rasterization.
func DecodeConfig(r io.Reader) (image.Config, error) {
	var cfg image.Config
	icon, err := oksvg.ReadIconStream(r)
	if err != nil {
		return cfg, err
	}
	if len(icon.SVGPaths) == 0 && icon.ViewBox.W == 0 && icon.ViewBox.H == 0 {
		return cfg, fmt.Errorf("invalid or empty SVG: no svg element or paths found")
	}

	w := int(math.Round(icon.ViewBox.W))
	h := int(math.Round(icon.ViewBox.H))
	if w <= 0 || h <= 0 {
		w, h = 512, 512
	}

	cfg.ColorModel = color.RGBAModel
	cfg.Width = w
	cfg.Height = h
	return cfg, nil
}

func init() {
	image.RegisterFormat("svg", "<?xml", Decode, DecodeConfig)
	image.RegisterFormat("svg", "<svg", Decode, DecodeConfig)
	image.RegisterFormat("svg", "\n<svg", Decode, DecodeConfig)
	image.RegisterFormat("svg", " <svg", Decode, DecodeConfig)
	image.RegisterFormat("svg", "\t<svg", Decode, DecodeConfig)
	image.RegisterFormat("svg", "\r\n<svg", Decode, DecodeConfig)
	image.RegisterFormat("svg", "<!--", Decode, DecodeConfig)
	image.RegisterFormat("svg", "<!DOCTYPE svg", Decode, DecodeConfig)
}
