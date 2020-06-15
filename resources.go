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

	// Write fonts
	if len(pdf.fonts) > 0 {
		pdf.outln("  /Font <<")
		for _, font := range pdf.fonts {
			pdf.outln(fmt.Sprintf("    /Font-%s %d 0 R", font.family, font.id))
		}
		pdf.outln("  >>")
	}

	// Write XObjects
	if len(pdf.tplObjIds)+len(pdf.images) > 0 {
		pdf.outln("  /XObject <<")
		for tplName, id := range pdf.tplObjIds {
			pdf.outln(fmt.Sprintf("    %s %d 0 R", tplName, id))
		}
		for _, image := range pdf.images {
			pdf.outln(fmt.Sprintf("    /Image%d %d 0 R", image.id, image.id))
		}
		pdf.outln("  >>")
	}

	// Write Shadings
	if len(pdf.shadings) > 0 {
		pdf.outln("  /Shading <<")
		for i := 0; i < len(pdf.shadings); i++ {
			shading := pdf.shadings[i]
			pdf.outln(fmt.Sprintf("    /Sh%d %d 0 R", 1, shading.id))
		}
		pdf.outln("  >>")
	}

	pdf.outln(">>")
	pdf.outln("endobj\n")
}
