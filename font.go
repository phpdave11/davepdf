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

func (pdf *Pdf) SetFontFamily(fontFamily string) {
	pdf.fontFamily = "/FONT1"
}

func (pdf *Pdf) SetFontSize(fontSize int) {
	pdf.fontSize = fontSize
}

func (pdf *Pdf) writeFonts() {
	pdf.newObj(pdf.font.id)
	pdf.outln("<<")
	pdf.outln("  /Type /Font")
	pdf.outln("  /Subtype /Type1")
	pdf.outln("  /Name /FONT1")
	pdf.outln("  /BaseFont /Times-Roman")
	pdf.outln(">>")
	pdf.outln("endobj\n")
}
