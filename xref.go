package davepdf

import "fmt"

type PdfXrefTable struct {
	id int
}

func (pdf *Pdf) writeXref() {
	pdf.newObjId()
	pdf.outln("xref")
	pdf.outln(fmt.Sprintf("0 %d", pdf.n-1))
	pdf.outln("0000000000 65535 f ")

	for i := 1; i < len(pdf.offsets); i++ {
		pdf.outln(fmt.Sprintf("%010d 00000 n ", pdf.offsets[i]))
	}
}

func (pdf *Pdf) writeTrailer() {
	pdf.outln("trailer")
	pdf.outln("<<")
	pdf.outln(fmt.Sprintf("  /Size %d", pdf.n-1))
	pdf.outln(fmt.Sprintf("  /Root %d 0 R", pdf.catalog.id))
	pdf.outln(">>")
	pdf.outln("startxref")
	pdf.outln(fmt.Sprintf("%d", pdf.offsets[pdf.n-1]))
	pdf.outln("%%EOF")
}
