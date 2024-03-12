package lens

import (
	"image"
	"math"

	"golang.org/x/image/math/fixed"
)

type PixelCoords struct {
	x      fixed.Int26_6
	y      fixed.Int26_6
	Bounds fixed.Rectangle26_6
}

type BoudingBox struct {
	Bounds      image.Rectangle
	CenterPerX  float64
	CenterPerY  float64
	PerWidth    float64
	PerHeight   float64
	PixelCoords PixelCoords
}

type Segment struct {
	Text       string
	BoudingBox BoudingBox
}

type SegmentOptions struct {
	regions []interface{}
	bounds  image.Rectangle
}

func NewSegment(text string, opts SegmentOptions) *Segment {
	return &Segment{
		Text:       text,
		BoudingBox: optsToBoundingBox(opts),
	}
}

func optsToBoundingBox(opts SegmentOptions) BoudingBox {
	var box = BoudingBox{
		Bounds:     opts.bounds,
		CenterPerX: opts.regions[0].(float64),
		CenterPerY: opts.regions[1].(float64),
		PerWidth:   opts.regions[2].(float64),
		PerHeight:  opts.regions[3].(float64),
	}
	box.PixelCoords = box.toPixelCoords()
	return box
}

func (box *BoudingBox) toPixelCoords() PixelCoords {
	imgWidth := float64(box.Bounds.Dx())
	imgHeight := float64(box.Bounds.Dy())
	width := box.PerWidth * imgWidth
	height := box.PerHeight * imgHeight
	x := (box.CenterPerX * imgWidth) - (width / 2)
	y := (box.CenterPerY * imgHeight) - (height / 2)
	return PixelCoords{
		x: float64ToFixed(x),
		y: float64ToFixed(y),
		Bounds: fixed.Rectangle26_6{
			Max: fixed.Point26_6{
				X: float64ToFixed(width),
				Y: float64ToFixed(height),
			},
			Min: fixed.Point26_6{
				X: fixed.I(0),
				Y: fixed.I(0),
			},
		},
	}
}

func float64ToFixed(value float64) fixed.Int26_6 {
	fixedValue := fixed.Int26_6(math.Round(value * 64))
	return fixedValue
}
