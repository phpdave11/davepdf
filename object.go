package davepdf

type PdfObject struct {
	id         int
	dictionary map[string]string
	data       []byte
}

func (pdf *Pdf) newObj(n int) {
	pdf.offsets[n] = pdf.w.Len()
	pdf.outln(fmt.Sprintf("%d 0 obj", n))
}

func (pdf *Pdf) newObjId() {
	pdf.n++
}
