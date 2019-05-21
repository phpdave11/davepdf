package davepdf

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
