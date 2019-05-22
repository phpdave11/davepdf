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

func (pdf *Pdf) SetFontFamily(fontFamily string) {
	pdf.fontFamily = "/FONT1"
}

func (pdf *Pdf) SetFontSize(fontSize int) {
	pdf.fontSize = fontSize
}

func (pdf *Pdf) SetFillColorCMYK(c, m, y, k int) {
	pdf.fillColor = &Color{colorSpace: ColorSpaceCMYK, cmyk: &CMYKColor{C: c, M: m, Y: y, K: k}}

	pdf.page.instructions.add(fmt.Sprintf("%.5f %.5f %.5f %.5f k", float64(pdf.fillColor.cmyk.C)/100, float64(pdf.fillColor.cmyk.M)/100, float64(pdf.fillColor.cmyk.Y)/100, float64(pdf.fillColor.cmyk.K)/100), "set fill color (cmyk)")
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
	pdf.newObj(pdf.catalog.id)
	pdf.outln("<<")
	pdf.outln("  /Type /Catalog")
	pdf.outln(fmt.Sprintf("  /Pages %d 0 R", pdf.pageTree.id))
	pdf.outln(">>")
	pdf.outln("endobj\n")

	// write page tree
	pdf.newObj(pdf.pageTree.id)
	pdf.outln("<<")
	pdf.outln("  /Type /Pages")
	pdf.outln("  /Count 1")
	pdf.outln(fmt.Sprintf("  /Kids [%d 0 R]", pdf.page.id))
	pdf.outln(">>")
	pdf.outln("endobj\n")

	// write resources
	pdf.newObj(pdf.resources.id)
	pdf.outln("<<")
	pdf.outln("  /ProcSet [/PDF /Text /ImageB /ImageC /ImageI]")
	pdf.outln("  /Font <<")
	pdf.outln(fmt.Sprintf("    /FONT1 %d 0 R", pdf.font.id))
	pdf.outln("  >>")
	pdf.outln("  /XObject <<")
	for tplName, id := range pdf.tplObjIds {
		pdf.outln(fmt.Sprintf("    %s %d 0 R", tplName, id))
	}
	pdf.outln("  >>")
	pdf.outln(">>")
	pdf.outln("endobj\n")

	// write fonts
	pdf.newObj(pdf.font.id)
	pdf.outln("<<")
	pdf.outln("  /Type /Font")
	pdf.outln("  /Subtype /Type1")
	pdf.outln("  /Name /FONT1")
	pdf.outln("  /BaseFont /Times-Roman")
	pdf.outln(">>")
	pdf.outln("endobj\n")

	// write page
	pdf.newObj(pdf.page.id)
	pdf.outln("<<")
	pdf.outln("  /Type /Page")
	pdf.outln("  /MediaBox [0 0 612 500]")
	//pdf.outln("  /MediaBox [0 0 595.28 841.89]")
	pdf.outln(fmt.Sprintf("  /Parent %d 0 R", pdf.pageTree.id))
	pdf.outln(fmt.Sprintf("  /Contents %d 0 R", pdf.page.contents.id))
	pdf.outln(fmt.Sprintf("  /Resources %d 0 R", pdf.resources.id))
	pdf.outln(">>")
	pdf.outln("endobj\n")

	// write page contents
	pdf.newObj(pdf.page.contents.id)
	pdf.outln("<<")
	pdf.outln(fmt.Sprintf("  /Length %d", len(pdf.page.contents.data)))
	pdf.outln(">>")
	pdf.outln("stream")
	pdf.outln(pdf.page.instructions.String())
	pdf.outln("endstream")
	pdf.outln("endobj\n")

	// write xref
	pdf.newObjId()
	pdf.outln("xref")
	pdf.outln(fmt.Sprintf("0 %d", pdf.n-1))
	pdf.outln("0000000000 65535 f ")

	for i := 1; i < len(pdf.offsets); i++ {
		pdf.outln(fmt.Sprintf("%010d 00000 n ", pdf.offsets[i]))
	}

	// write trailer
	pdf.outln("trailer")
	pdf.outln("<<")
	pdf.outln(fmt.Sprintf("  /Size %d", pdf.n-1))
	pdf.outln(fmt.Sprintf("  /Root %d 0 R", pdf.catalog.id))
	pdf.outln(">>")
	pdf.outln("startxref")
	pdf.outln(fmt.Sprintf("%d", pdf.offsets[pdf.n-1]))
	pdf.outln("%%EOF")

	fmt.Print(string(pdf.w.Bytes()))
}
