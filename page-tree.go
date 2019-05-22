package davepdf

import "fmt"

type PdfPageTree struct {
	id    int
	pages []*PdfPage
}

func (pdf *Pdf) newPageTree() *PdfPageTree {
	pageTree := &PdfPageTree{}

	pdf.newObjId()
	pageTree.id = pdf.n

	return pageTree
}

func (pdf *Pdf) writePageTree() {
	pdf.newObj(pdf.pageTree.id)
	pdf.outln("<<")
	pdf.outln("  /Type /Pages")
	pdf.outln("  /Count 1")
	pdf.outln(fmt.Sprintf("  /Kids [%d 0 R]", pdf.page.id))
	pdf.outln(">>")
	pdf.outln("endobj\n")
}
