package davepdf

type PdfObject struct {
	id         int
	dictionary map[string]string
	data       []byte
}

