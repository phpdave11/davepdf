package davepdf

import (
	"bytes"
	"fmt"
	"github.com/phpdave11/gofpdi"
	"math"
)

type Pdf struct {
	n          int
	w          *bytes.Buffer
	catalog    *PdfCatalog
	resources  *PdfResources
	font       *PdfFont // temporary while we only use 1 font
	pageTree   *PdfPageTree
	xref       *PdfXrefTable
	fpdi       *gofpdi.Importer
	tplObjIds  map[string]int
	offsets    map[int]int
	objects    []*PdfObject
	fontFamily string
	fontSize   int
	x          float64
	y          float64
	k          float64
	h          float64
	page       *PdfPage
	fillColor  *Color
}

func NewPdf() *Pdf {
	pdf := &Pdf{}
	pdf.offsets = make(map[int]int, 0)
	pdf.w = new(bytes.Buffer)
	pdf.catalog = pdf.newCatalog()
	pdf.pageTree = pdf.newPageTree()
	pdf.resources = pdf.newResources()
	pdf.font = pdf.newFont()
	pdf.catalog.pageTree = pdf.pageTree
	pdf.fpdi = gofpdi.NewImporter()
	pdf.tplObjIds = make(map[string]int, 0)

	pdf.k = 1.0
	pdf.h = 841.89

	pdf.outln("%PDF-1.4")
	pdf.outln("%ABCD\n")

	return pdf
}

func (pdf *Pdf) SetXY(x, y float64) {
	pdf.x = x
	pdf.y = y
}

func (pdf *Pdf) Text(text string) {
	pdf.page.instructions.add("BT", "begin text")
	pdf.page.instructions.add(fmt.Sprintf("  %s %d Tf", pdf.fontFamily, pdf.fontSize), "set font family and font size")
	pdf.page.instructions.add(fmt.Sprintf("  %f %f Td", pdf.x, pdf.y), "set position to draw text")
	pdf.page.instructions.add(fmt.Sprintf("  (%s)Tj", text), "write text")
	pdf.page.instructions.add("ET", "end text")
}

func (pdf *Pdf) Rect(x, y, w, h float64, mode string) {
	pdf.page.instructions.add(fmt.Sprintf("%.5f %.5f %.5f %.5f re %s", x, y, w, h, mode), "draw rectangle")
}

func (pdf *Pdf) Circle(x, y, r float64, style string) {
	pdf.Ellipse(x, y, r, r, style)
}

func (pdf *Pdf) Ellipse(x, y, rx, ry float64, style string) {
	var op string
	var lx, ly, k /*, h*/ float64

	if style == "F" {
		op = "f"
	} else if style == "FD" || style == "DF" {
		op = "B"
	} else {
		op = "S"
	}

	M_SQRT2 := math.Sqrt(2)

	lx = 4.0 / 3.0 * (M_SQRT2 - 1.0) * rx
	ly = 4.0 / 3.0 * (M_SQRT2 - 1.0) * ry
	k = pdf.k

	pdf.page.instructions.add(fmt.Sprintf("%.2F %.2F m %.2F %.2F %.2F %.2F %.2F %.2F c", (x+rx)*k, (y)*k, (x+rx)*k, (y-ly)*k, (x+lx)*k, (y-ry)*k, x*k, (y-ry)*k), "draw ellipse part 1")
	pdf.page.instructions.add(fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", (x-lx)*k, (y-ry)*k, (x-rx)*k, (y-ly)*k, (x-rx)*k, (y)*k), "draw ellipse part 2")
	pdf.page.instructions.add(fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", (x-rx)*k, (y+ly)*k, (x-lx)*k, (y+ry)*k, x*k, (y+ry)*k), "draw ellipse part 3")
	pdf.page.instructions.add(fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c %s", (x+lx)*k, (y+ry)*k, (x+rx)*k, (y+ly)*k, (x+rx)*k, (y)*k, op), "draw ellipse part 4")

	// the following code makes everything start from the top based on height of page (pdf.h)
	/*
		h = pdf.h
		instructions += fmt.Sprintf("%-60s %% draw ellipse part 1\n", fmt.Sprintf("%.2F %.2F m %.2F %.2F %.2F %.2F %.2F %.2F c", (x+rx)*k, (h-y)*k, (x+rx)*k, (h-(y-ly))*k, (x+lx)*k, (h-(y-ry))*k, x*k, (h-(y-ry))*k))
		instructions += fmt.Sprintf("%-60s %% draw ellipse part 2\n", fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", (x-lx)*k, (h-(y-ry))*k, (x-rx)*k, (h-(y-ly))*k, (x-rx)*k, (h-y)*k))
		instructions += fmt.Sprintf("%-60s %% draw ellipse part 3\n", fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", (x-rx)*k, (h-(y+ly))*k, (x-lx)*k, (h-(y+ry))*k, x*k, (h-(y+ry))*k))
		instructions += fmt.Sprintf("%-60s %% draw ellipse part 4\n", fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c %s", (x+lx)*k, (h-(y+ry))*k, (x+rx)*k, (h-(y+ly))*k, (x+rx)*k, (h-y)*k, op))
	*/
}

func (pdf *Pdf) AddPage() *PdfPage {
	page := pdf.newPage()

	pdf.page = page

	return page
}

func (pdf *Pdf) outln(s string) {
	pdf.w.WriteString(s)
	pdf.w.WriteString("\n")
}

func (pdf *Pdf) out(s string) {
	pdf.w.WriteString(s)
}

func (pdf *Pdf) newObj(n int) {
	pdf.offsets[n] = pdf.w.Len()
	pdf.outln(fmt.Sprintf("%d 0 obj", n))
}

func (pdf *Pdf) newObjId() {
	pdf.n++
}

func (pdf *Pdf) Write() {
	// write catalog
	pdf.writeCatalog()

	// write page tree
	pdf.writePageTree()

	// write resources
	pdf.writeResources()

	// write fonts
	pdf.writeFonts()

	// write page
	pdf.writePage()

	// write page contents
	pdf.writeContents()

	// write xref
	pdf.writeXref()

	// write trailer
	pdf.writeTrailer()

	// output PDF
	fmt.Print(string(pdf.w.Bytes()))
}
