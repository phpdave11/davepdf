package davepdf

import "fmt"

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

func (pdf *Pdf) writeCatalog() {
	pdf.newObj(pdf.catalog.id)
	pdf.outln("<<")
	pdf.outln("  /Type /Catalog")
	pdf.outln(fmt.Sprintf("  /Pages %d 0 R", pdf.pageTree.id))
	pdf.outln(">>")
	pdf.outln("endobj\n")
}
