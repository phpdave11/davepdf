package davepdf

import (
	"fmt"
	"strings"
)

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

func (pdf *Pdf) AddPage() *PdfPage {
	page := pdf.newPage()

	pdf.page = page

	return page
}

func (pdf *Pdf) writePageTree() {
	// Kids (child pages of page tree) - e.g. /Kids [3 0 R 21 0 R]
	kids := ""
	for _, page := range pdf.pageTree.pages {
		kids += fmt.Sprintf("%d 0 R ", page.id)
	}
	// Trim leading space
	kids = strings.TrimSpace(kids)

	pdf.newObj(pdf.pageTree.id)
	pdf.outln("<<")
	pdf.outln("  /Type /Pages")
	pdf.outln(fmt.Sprintf("  /Count %d", len(pdf.pageTree.pages)))
	pdf.outln(fmt.Sprintf("  /Kids [%s]", kids))
	pdf.outln(">>")
	pdf.outln("endobj\n")
}
