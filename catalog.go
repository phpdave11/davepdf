package davepdf

type PdfCatalog struct {
	id       int
	pageTree *PdfPageTree
}

func (pdf *Pdf) newCatalog() *PdfCatalog {
	catalog := &PdfCatalog{}

	pdf.newObjId()
	catalog.id = pdf.n

	return catalog
}
