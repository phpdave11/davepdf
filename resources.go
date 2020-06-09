package davepdf

import "fmt"

type PdfResources struct {
	id      int
	objects []*PdfObject
}

func (pdf *Pdf) newResources() *PdfResources {
	resources := &PdfResources{}

	pdf.newObjId()
	resources.id = pdf.n

	return resources
}

func (pdf *Pdf) writeResources() {
	pdf.newObj(pdf.resources.id)
	pdf.outln("<<")
	pdf.outln("  /ProcSet [/PDF /Text /ImageB /ImageC /ImageI]")
	pdf.outln("  /Font <<")
	for _, font := range pdf.fonts {
		pdf.outln(fmt.Sprintf("    /FONT%d %d 0 R", font.id, font.id))
	}
	pdf.outln("  >>")
	pdf.outln("  /XObject <<")
	for tplName, id := range pdf.tplObjIds {
		pdf.outln(fmt.Sprintf("    %s %d 0 R", tplName, id))
	}
	pdf.outln("  >>")
	pdf.outln("  /Shading <<")
	for i := 0; i < len(pdf.shadings); i++ {
		shading := pdf.shadings[i]
		//pdf.outln(fmt.Sprintf("    /Sh%d %d 0 R", shading.id, shading.id))
		pdf.outln(fmt.Sprintf("    /Sh%d %d 0 R", 1, shading.id))
	}
	pdf.outln("  >>")
	pdf.outln(">>")
	pdf.outln("endobj\n")
}
