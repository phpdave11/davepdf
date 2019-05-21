package davepdf

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

