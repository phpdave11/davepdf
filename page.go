package davepdf

type PdfPage struct {
	id           int
	parent       *PdfPageTree
	instructions *PdfInstructions // to be compiled into contents below
	contents     *PdfContents
}

func (pdf *Pdf) newPage() *PdfPage {
	page := &PdfPage{}

	pdf.newObjId()
	page.id = pdf.n

	page.parent = pdf.pageTree
	page.contents = pdf.newContents()
	page.instructions = pdf.newInstructions()

	pdf.pageTree.pages = append(pdf.pageTree.pages, page)

	return page
}
