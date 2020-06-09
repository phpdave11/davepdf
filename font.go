package davepdf

import "strconv"

//import "github.com/davecgh/go-spew/spew"

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
	found := false

	for _, font := range pdf.fonts {
		if font.family == fontFamily {
			pdf.fontFamily = "/FONT" + strconv.Itoa(font.id)
			found = true
			break
		}
	}

	if !found {
		font := pdf.newFont()
		font.family = fontFamily
		pdf.fontFamily = "/FONT" + strconv.Itoa(font.id)
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
		pdf.outln("  /Name /FONT" + strconv.Itoa(font.id))
		pdf.outln("  /BaseFont /" + font.family)
		pdf.outln(">>")
		pdf.outln("endobj\n")
	}
}
