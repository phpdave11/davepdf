package davepdf

type PdfPage struct {
	id       int
	parent   *PdfPageTree
	contents *PdfContents
}

func (pdf *Pdf) newPage() *PdfPage {
	page := &PdfPage{}

	pdf.newObjId()
	page.id = pdf.n

	page.parent = pdf.pageTree
	page.contents = pdf.newContents()

	pdf.pageTree.pages = append(pdf.pageTree.pages, page)

	return page
}
