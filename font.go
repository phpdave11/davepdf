package davepdf

type PdfFont struct {
	id     int
	family string
}

func (pdf *Pdf) newFont() *PdfFont {
	font := &PdfFont{}

	pdf.newObjId()
	font.id = pdf.n

	pdf.fonts = append(pdf.fonts, font)

	return font
}

func (pdf *Pdf) SetFontFamily(fontFamily string) {
	validFonts := []string{
		"Times-Roman",
		"Times-Bold",
		"Times-Italic",
		"Times-BoldItalic",
		"Courier",
		"Courier-Bold",
		"Courier-Oblique",
		"Courier-BoldOblique",
		"Helvetica",
		"Helvetica-Bold",
		"Helvetica-Oblique",
		"Helvetica-BoldOblique",
		"Symbol",
		"ZapfDingbats"}

	valid := false
	for _, font := range validFonts {
		if font == fontFamily {
			valid = true
			break
		}
	}

	if !valid {
		panic("Invalid font: " + fontFamily)
	}

	found := false

	for _, font := range pdf.fonts {
		if font.family == fontFamily {
			pdf.fontFamily = "/Font-" + font.family
			found = true
			break
		}
	}

	if !found {
		font := pdf.newFont()
		font.family = fontFamily
		pdf.fontFamily = "/Font-" + font.family
	}
}

func (pdf *Pdf) SetFontSize(fontSize int) {
	pdf.fontSize = fontSize
}

func (pdf *Pdf) writeFonts() {
	for _, font := range pdf.fonts {
		pdf.newObj(font.id)
		pdf.outln("<<")
		pdf.outln("  /Type /Font")
		pdf.outln("  /Subtype /Type1")
		pdf.outln("  /Name /Font-" + font.family)
		pdf.outln("  /BaseFont /" + font.family)
		pdf.outln(">>")
		pdf.outln("endobj\n")
	}
}
