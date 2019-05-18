package davepdf

import (
	"bytes"
	"fmt"
	"math"
	//"github.com/phpdave11/gofpdi"
)

type Pdf struct {
	n          int
	w          *bytes.Buffer
	catalog    *PdfCatalog
	pageTree   *PdfPageTree
	xref       *PdfXrefTable
	offsets    map[int]int
	objects    []*PdfObject
	fontFamily string
	fontSize   int
	x          float64
	y          float64
	k          float64
	h          float64
	page       *PdfPage
	fillColor  *CMYKColor
}

type CMYKColor struct {
	C int
	Y int
	M int
	K int
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
	pdf.k = 1.0
	pdf.h = 841.89

	return pdf
}

func (pdf *Pdf) SetFontFamily(fontFamily string) {
	pdf.fontFamily = "/FONT1"
}

func (pdf *Pdf) SetFontSize(fontSize int) {
	pdf.fontSize = fontSize
}

func (pdf *Pdf) SetFillColor(c *CMYKColor) {
	pdf.fillColor = c

	var instructions string
	instructions += fmt.Sprintf("%-60s %% set fill color (cmyk)\n", fmt.Sprintf("%.5f %.5f %.5f %.5f k", float64(c.C)/100, float64(c.M)/100, float64(c.Y)/100, float64(c.K)/100))
	pdf.page.contents.data = append(pdf.page.contents.data, []byte(instructions)...)
}

func (pdf *Pdf) SetXY(x, y float64) {
	pdf.x = x
	pdf.y = y
}

func (pdf *Pdf) Text(text string) {
	var instructions string

	instructions += fmt.Sprintf("%-60s %% begin text\n", "BT")
	instructions += fmt.Sprintf("  %-58s %% set font family and font size\n", fmt.Sprintf("%s %d Tf", pdf.fontFamily, pdf.fontSize))
	instructions += fmt.Sprintf("  %-58s %% set position to draw text\n", fmt.Sprintf("%f %f Td", pdf.x, pdf.y))
	instructions += fmt.Sprintf("  %-58s %% write text\n", fmt.Sprintf("(%s)Tj", text))
	instructions += fmt.Sprintf("%-60s %% end text\n", "ET")

	pdf.page.contents.data = append(pdf.page.contents.data, []byte(instructions)...)
}

func (pdf *Pdf) Rect(x, y, w, h float64, mode string) {
	var instructions string

	instructions += fmt.Sprintf("%-60s %% draw rectangle\n", fmt.Sprintf("%.5f %.5f %.5f %.5f re %s", x, y, w, h, mode))

	pdf.page.contents.data = append(pdf.page.contents.data, []byte(instructions)...)
}

func (pdf *Pdf) Circle(x, y, r float64, style string) {
	pdf.Ellipse(x, y, r, r, style)
}

func (pdf *Pdf) Ellipse(x, y, rx, ry float64, style string) {
    var instructions string

	var op string
	var lx, ly, k/*, h*/ float64

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

	instructions += fmt.Sprintf("%-60s %% draw ellipse part 1\n", fmt.Sprintf("%.2F %.2F m %.2F %.2F %.2F %.2F %.2F %.2F c", (x+rx)*k, (y)*k, (x+rx)*k, ((y-ly))*k, (x+lx)*k, ((y-ry))*k, x*k, ((y-ry))*k))
	instructions += fmt.Sprintf("%-60s %% draw ellipse part 2\n", fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", (x-lx)*k, ((y-ry))*k, (x-rx)*k, ((y-ly))*k, (x-rx)*k, (y)*k))
	instructions += fmt.Sprintf("%-60s %% draw ellipse part 3\n", fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", (x-rx)*k, ((y+ly))*k, (x-lx)*k, ((y+ry))*k, x*k, ((y+ry))*k))
	instructions += fmt.Sprintf("%-60s %% draw ellipse part 4\n", fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c %s", (x+lx)*k, ((y+ry))*k, (x+rx)*k, ((y+ly))*k, (x+rx)*k, (y)*k, op))

	// the following code makes everything start from the top based on height of page (pdf.h)
	/*
	h = pdf.h
	instructions += fmt.Sprintf("%-60s %% draw ellipse part 1\n", fmt.Sprintf("%.2F %.2F m %.2F %.2F %.2F %.2F %.2F %.2F c", (x+rx)*k, (h-y)*k, (x+rx)*k, (h-(y-ly))*k, (x+lx)*k, (h-(y-ry))*k, x*k, (h-(y-ry))*k))
	instructions += fmt.Sprintf("%-60s %% draw ellipse part 2\n", fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", (x-lx)*k, (h-(y-ry))*k, (x-rx)*k, (h-(y-ly))*k, (x-rx)*k, (h-y)*k))
	instructions += fmt.Sprintf("%-60s %% draw ellipse part 3\n", fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", (x-rx)*k, (h-(y+ly))*k, (x-lx)*k, (h-(y+ry))*k, x*k, (h-(y+ry))*k))
	instructions += fmt.Sprintf("%-60s %% draw ellipse part 4\n", fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c %s", (x+lx)*k, (h-(y+ry))*k, (x+rx)*k, (h-(y+ly))*k, (x+rx)*k, (h-y)*k, op))
    */

	pdf.page.contents.data = append(pdf.page.contents.data, []byte(instructions)...)
}

func (pdf *Pdf) AddPage() *PdfPage {
	page := &PdfPage{}
	page.parent = pdf.pageTree
	page.contents = &PdfContents{}
	pdf.page = page

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
	pdf.outln("  /XObject <<")
	//pdf.outln("    /GOFPDITPL0 7 0 R")
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
	pdf.outln("  /BaseFont /Times-Roman")
	pdf.outln(">>")
	pdf.outln("endobj\n")
	/*
		// init gofpdi
		fpdi := gofpdi.NewImporter()
		fpdi.SetSourceFile("/Users/dave/Desktop/PDFPL110.pdf")
		fpdi.SetNextObjectID(7)
		tpl1 := fpdi.ImportPage(1, "/MediaBox")

		// get gofpdi template values
		tplName, scaleX, scaleY, tX, tY := fpdi.UseTemplate(tpl1, 0, -400, 400, 0)
		str := fmt.Sprintf("q 0 J 1 w 0 j 0 G 0 g q %.4F 0 0 %.4F %.4F %.4F cm %s Do Q Q", scaleX, scaleY, tX, tY, tplName)
	*/

	// add page
	page := pdf.AddPage()
	//page.contents.data = []byte("BT /FONT1 18 Tf 0 0 Td (Hello World) Tj ET " + str)

	pdf.SetFontFamily("Times-Roman")
	pdf.SetFontSize(18)
	pdf.SetXY(10, 600)
	pdf.SetXY(0, 0)
	pdf.Text( /*"こんにちは and " + */ "Hello World!")

	//pdf.SetFillColor(&CMYKColor{C: 48, M: 32, Y: 0, K: 0})
	pdf.SetFillColor(&CMYKColor{C: 0, M: 81, Y: 81, K: 45})
	pdf.Rect(10, 200, 250, 50, "F")

	pdf.SetFillColor(&CMYKColor{C: 26, M: 0, Y: 99, K: 13})
	//pdf.Ellipse(100, 50, 30, 20, "D")
	pdf.Circle(110, 300, 70, "F")

	// write page
	pdf.newObj()
	pdf.outln(fmt.Sprintf("%d 0 obj", pdf.n))
	pdf.outln("<<")
	pdf.outln("  /Type /Page")
	pdf.outln("  /MediaBox [0 0 612 500]")
	//pdf.outln("  /MediaBox [0 0 595.28 841.89]")
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

	/*
		// write imported objects
		fpdi.PutFormXobjects()

		objs := fpdi.GetImportedObjects()
		for i := 7; i < len(objs)+7; i++ {
			pdf.newObj()
			pdf.outln(fmt.Sprintf("%d 0 obj", i))
			pdf.outln(objs[i])
		}
	*/

	// write xref
	pdf.newObj()
	pdf.outln("xref")
	pdf.outln(fmt.Sprintf("0 %d", pdf.n))
	pdf.outln("0000000000 65535 f ")
	for i := 0; i < len(pdf.offsets)-1; i++ {
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

	fmt.Print(string(pdf.w.Bytes()))
}
