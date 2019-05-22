package davepdf

import "fmt"

type PdfPage struct {
	id           int
	parent       *PdfPageTree
	instructions *PdfInstructions // to be compiled into contents below
	contents     *PdfContents
}

func (pdf *Pdf) newPage() *PdfPage {
	page := &PdfPage{}

	pdf.newObjId()
	page.id = pdf.n

	page.parent = pdf.pageTree
	page.contents = pdf.newContents()
	page.instructions = pdf.newInstructions()

	pdf.pageTree.pages = append(pdf.pageTree.pages, page)

	return page
}

func (pdf *Pdf) writePage() {
	pdf.newObj(pdf.page.id)
	pdf.outln("<<")
	pdf.outln("  /Type /Page")
	pdf.outln("  /MediaBox [0 0 612 500]")
	pdf.outln(fmt.Sprintf("  /Parent %d 0 R", pdf.pageTree.id))
	pdf.outln(fmt.Sprintf("  /Contents %d 0 R", pdf.page.contents.id))
	pdf.outln(fmt.Sprintf("  /Resources %d 0 R", pdf.resources.id))
	pdf.outln(">>")
	pdf.outln("endobj\n")
}
