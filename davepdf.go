package davepdf

import (
	"bytes"
	"fmt"
)

type Pdf struct {
	n        int
	w        *bytes.Buffer
	catalog  *PdfCatalog
	pageTree *PdfPageTree
	xref     *PdfXrefTable
	offsets  map[int]int
	objects  []*PdfObject
}

type PdfCatalog struct {
	id       int
	pageTree *PdfPageTree
}

type PdfObject struct {
	id         int
	dictionary map[string]string
	data       []byte
}

type PdfXrefTable struct {
	id int
}

type PdfPageTree struct {
	id    int
	pages []*PdfPage
}

type PdfPage struct {
	id        int
	parent    *PdfPageTree
	contents  *PdfContents
	resourecs *PdfResources
}

type PdfContents struct {
	id   int
	data []byte
}

type PdfResources struct {
	id      int
	objects []*PdfObject
}

func NewPdf() *Pdf {
	pdf := &Pdf{}
	pdf.offsets = make(map[int]int, 0)
	pdf.w = new(bytes.Buffer)
	pdf.pageTree = &PdfPageTree{}
	pdf.catalog = &PdfCatalog{}
	pdf.catalog.pageTree = pdf.pageTree

	return pdf
}

func (pdf *Pdf) AddPage() *PdfPage {
	page := &PdfPage{}
	page.parent = pdf.pageTree
	page.contents = &PdfContents{}

	pdf.pageTree.pages = append(pdf.pageTree.pages, page)

	return page
}

func (pdf *Pdf) outln(s string) {
	pdf.w.WriteString(s)
	pdf.w.WriteString("\n")
}

func (pdf *Pdf) out(s string) {
	pdf.w.WriteString(s)
}

func (pdf *Pdf) newObj() {
	pdf.offsets[pdf.n] = pdf.w.Len()
	pdf.n++
}

func (pdf *Pdf) Write() {
	pdf.outln("%PDF-1.4")
	pdf.outln("%ABCD\n")

	// write catalog
	pdf.newObj()
	pdf.outln(fmt.Sprintf("%d 0 obj", pdf.n))
	pdf.outln("<<")
	pdf.outln("  /Type /Catalog")
	pdf.outln("  /Pages 2 0 R")
	pdf.outln(">>")
	pdf.outln("endobj\n")

	// write page tree
	pdf.newObj()
	pdf.outln(fmt.Sprintf("%d 0 obj", pdf.n))
	pdf.outln("<<")
	pdf.outln("  /Type /Pages")
	pdf.outln("  /Count 1")
	pdf.outln("  /Kids [5 0 R]")
	pdf.outln(">>")
	pdf.outln("endobj\n")

	// write resources
	pdf.newObj()
	pdf.outln(fmt.Sprintf("%d 0 obj", pdf.n))
	pdf.outln("<<")
	pdf.outln("  /ProcSet [/PDF /Text /ImageB /ImageC /ImageI]")
	pdf.outln("  /Font <<")
	pdf.outln("    /FONT1 4 0 R")
	pdf.outln("  >>")
	pdf.outln(">>")
	pdf.outln("endobj\n")

	// write fonts
	pdf.newObj()
	pdf.outln(fmt.Sprintf("%d 0 obj", pdf.n))
	pdf.outln("<<")
	pdf.outln("  /Type /Font")
	pdf.outln("  /Subtype /Type1")
	pdf.outln("  /Name /FONT1")
	pdf.outln("  /BaseFont /Helvetica")
	pdf.outln("  /Encoding /StandardEncoding")
	pdf.outln(">>")
	pdf.outln("endobj\n")

	// add page
	page := pdf.AddPage()
	page.contents.data = []byte("BT /FONT1 18 Tf 0 0 Td (Hello World) Tj ET")

	// write page
	pdf.newObj()
	pdf.outln(fmt.Sprintf("%d 0 obj", pdf.n))
	pdf.outln("<<")
	pdf.outln("  /Type /Page")
	pdf.outln("  /MediaBox [0 0 612 792]")
	pdf.outln("  /Parent 2 0 R")
	pdf.outln("  /Contents 6 0 R")
	pdf.outln("  /Resources 3 0 R")
	pdf.outln(">>")
	pdf.outln("endobj\n")

	// write page contents
	pdf.newObj()
	pdf.outln(fmt.Sprintf("%d 0 obj", pdf.n))
	pdf.outln("<<")
	pdf.outln(fmt.Sprintf("  /Length %d", len(page.contents.data)))
	pdf.outln(">>")
	pdf.outln("stream")
	pdf.outln(string(page.contents.data))
	pdf.outln("endstream")
	pdf.outln("endobj\n")

	// write xref
	pdf.newObj()
	pdf.outln("xref")
	pdf.outln(fmt.Sprintf("0 %d", pdf.n))
	pdf.outln("0000000000 65535 f ")
	for i := 0; i < len(pdf.offsets); i++ {
		pdf.outln(fmt.Sprintf("%010d 00000 n ", pdf.offsets[i]))
	}

	// write trailer
	pdf.outln("trailer")
	pdf.outln("<<")
	pdf.outln(fmt.Sprintf("  /Size %d", pdf.n))
	pdf.outln("  /Root 1 0 R")
	pdf.outln(">>")
	pdf.outln("startxref")
	pdf.outln(fmt.Sprintf("%d", pdf.offsets[pdf.n-1]))
	pdf.outln("%%EOF")

	fmt.Println(string(pdf.w.Bytes()))
}
