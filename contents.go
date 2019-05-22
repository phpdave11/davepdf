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
	pdf.newObj(pdf.page.contents.id)
	pdf.outln("<<")
	pdf.outln(fmt.Sprintf("  /Length %d", len(pdf.page.contents.data)))
	pdf.outln(">>")
	pdf.outln("stream")
	pdf.outln(pdf.page.instructions.String())
	pdf.outln("endstream")
	pdf.outln("endobj\n")
}
