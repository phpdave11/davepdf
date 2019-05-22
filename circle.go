package davepdf

func (pdf *Pdf) Circle(x, y, r float64, style string) {
	pdf.Ellipse(x, y, r, r, style)
}
