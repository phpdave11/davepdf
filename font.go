package davepdf

type PdfFont struct {
	id int
}

func (pdf *Pdf) newFont() *PdfFont {
	font := &PdfFont{}

	pdf.newObjId()
	font.id = pdf.n

	return font
}

