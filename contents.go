package davepdf

import "fmt"

type PdfContents struct {
	id   int
	data []byte
}

func (pdf *Pdf) newContents() *PdfContents {
	contents := &PdfContents{}

	pdf.newObjId()
	contents.id = pdf.n

	return contents
}

func (pdf *Pdf) writeContents() {
	for _, page := range pdf.pageTree.pages {
		inst := page.instructions.String()
		pdf.newObj(page.contents.id)
		pdf.outln("<<")
		pdf.outln(fmt.Sprintf("  /Length %d", len(inst)))
		pdf.outln(">>")
		pdf.outln("stream")
		pdf.outln(inst)
		pdf.outln("endstream")
		pdf.outln("endobj\n")
	}
}
