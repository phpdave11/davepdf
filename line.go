package davepdf

import (
	"fmt"
)

// Line width
func (pdf *Pdf) SetLineWidth(w float64) {
	pdf.lineWidth = w
	pdf.page.instructions.add(fmt.Sprintf("%f w", pdf.lineWidth), "set line width")
}

// Line cap style
// With type assertions in Go 1.9, you can simply add = where you define the type.
// This tells the compiler that LineCapStyle is an alternate name for uint.
type LineCapStyle = uint

const (
	LineCapStyleButt LineCapStyle = iota
	LineCapStyleRound
	LineCapStyleProjectingSquare
)

type LineCap struct {
	style LineCapStyle
}

func (pdf *Pdf) SetLineCapStyle(s uint) {
	pdf.lineCapStyle = &LineCap{style: s}
	pdf.page.instructions.add(fmt.Sprintf("%d J", pdf.lineCapStyle.style), "set line cap style")
}

// Line join style
type LineJoinStyle = uint

const (
	LineJoinStyleMiter LineJoinStyle = iota
	LineJoinStyleRound
	LineJoinStyleBevel
)

type LineJoin struct {
	style LineJoinStyle
}

func (pdf *Pdf) SetLineJoinStyle(s uint) {
	pdf.lineJoinStyle = &LineJoin{style: s}
	pdf.page.instructions.add(fmt.Sprintf("%d j", pdf.lineJoinStyle.style), "set line join style")
}

// Line dash pattern
type LineDashPattern = uint

const (
	LineDashPatternMiter LineDashPattern = iota
	LineDashPatternRound
	LineDashPatternBevel
)

type LineDash struct {
	array []uint
	phase uint
}

func (pdf *Pdf) SetLineDashPattern(array []uint, phase uint) {
	pdf.lineDashPattern = &LineDash{array: array, phase: phase}
	pdf.page.instructions.add(fmt.Sprintf("%d %d d", pdf.lineDashPattern.array, pdf.lineDashPattern.phase), "set line dash pattern")
}
