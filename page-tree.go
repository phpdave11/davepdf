package davepdf

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
