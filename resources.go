package davepdf

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
