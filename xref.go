package davepdf

import "fmt"

type PdfXrefTable struct {
	offset int
}

func (pdf *Pdf) newXrefTable() *PdfXrefTable {
	xrefTable := &PdfXrefTable{}
	return xrefTable
}

func (pdf *Pdf) writeXref() {
	pdf.newObjId()
	pdf.xref.offset = pdf.w.Len()
	pdf.outln("xref")
	pdf.outln(fmt.Sprintf("0 %d", pdf.n))
	pdf.outln("0000000000 65535 f ")

	for i := 1; i <= len(pdf.offsets); i++ {
		pdf.outln(fmt.Sprintf("%010d 00000 n ", pdf.offsets[i]))
	}
}

func (pdf *Pdf) writeTrailer() {
	pdf.outln("\ntrailer")
	pdf.outln("<<")
	pdf.outln(fmt.Sprintf("  /Size %d", pdf.n))
	pdf.outln(fmt.Sprintf("  /Root %d 0 R", pdf.catalog.id))
	pdf.outln(">>")
	pdf.outln("startxref")
	pdf.outln(fmt.Sprintf("%d", pdf.xref.offset))
	pdf.outln("%%EOF")
}
