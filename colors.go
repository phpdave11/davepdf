package davepdf

import "fmt"

type CMYKColor struct {
	C int
	Y int
	M int
	K int
}

type RGBColor struct {
	R int
	G int
	B int
}

type ColorSpace uint

const (
	ColorSpaceCMYK ColorSpace = iota
	ColorSpaceRGB
)

type Color struct {
	colorSpace ColorSpace
	cmyk       *CMYKColor
	rgb        *RGBColor
}

func (pdf *Pdf) SetFillColorCMYK(c, m, y, k int) {
	pdf.fillColor = &Color{colorSpace: ColorSpaceCMYK, cmyk: &CMYKColor{C: c, M: m, Y: y, K: k}}

	pdf.page.instructions.add(fmt.Sprintf("%.5f %.5f %.5f %.5f k", float64(pdf.fillColor.cmyk.C)/100, float64(pdf.fillColor.cmyk.M)/100, float64(pdf.fillColor.cmyk.Y)/100, float64(pdf.fillColor.cmyk.K)/100), "set fill color (cmyk)")
}

func (pdf *Pdf) SetStrokeColorCMYK(c, m, y, k int) {
	pdf.strokeColor = &Color{colorSpace: ColorSpaceCMYK, cmyk: &CMYKColor{C: c, M: m, Y: y, K: k}}

	pdf.page.instructions.add(fmt.Sprintf("%.5f %.5f %.5f %.5f K", float64(pdf.strokeColor.cmyk.C)/100, float64(pdf.strokeColor.cmyk.M)/100, float64(pdf.strokeColor.cmyk.Y)/100, float64(pdf.strokeColor.cmyk.K)/100), "set stroke color (cmyk)")
}
