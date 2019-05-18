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

type PdfXObject struct {
	id int
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
	id       int
	parent   *PdfPageTree
	contents *PdfContents
}

type PdfContents struct {
	id   int
	data []byte
}

type PdfResources struct {
	id      int
	objects []*PdfObject
}

type PdfFont struct {
	id int
}

func (pdf *Pdf) newPageTree() *PdfPageTree {
	pageTree := &PdfPageTree{}

	pdf.newObjId()
	pageTree.id = pdf.n

	return pageTree
}

func (pdf *Pdf) newCatalog() *PdfCatalog {
	catalog := &PdfCatalog{}

	pdf.newObjId()
	catalog.id = pdf.n

	return catalog
}

func (pdf *Pdf) newResources() *PdfResources {
	resources := &PdfResources{}

	pdf.newObjId()
	resources.id = pdf.n

	return resources
}

func (pdf *Pdf) newContents() *PdfContents {
	contents := &PdfContents{}

	pdf.newObjId()
	contents.id = pdf.n

	return contents
}

func (pdf *Pdf) newFont() *PdfFont {
	font := &PdfFont{}

	pdf.newObjId()
	font.id = pdf.n

	return font
}

func (pdf *Pdf) newPage() *PdfPage {
	page := &PdfPage{}

	pdf.newObjId()
	page.id = pdf.n

	page.parent = pdf.pageTree
	page.contents = pdf.newContents()

	pdf.pageTree.pages = append(pdf.pageTree.pages, page)

	return page
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

func (pdf *Pdf) UseImportedTemplate(tplid int, x, y, w, h float64) {
	var instructions string

	// get gofpdi template values
	tplName, scaleX, scaleY, tX, tY := pdf.fpdi.UseTemplate(tplid, x, y, w, h)
	instructions += fmt.Sprintf("q 0 J 1 w 0 j 0 G 0 g q %.4F 0 0 %.4F %.4F %.4F cm %s Do Q Q %% draw template\n", scaleX, scaleY, tX, tY, tplName)

	pdf.page.contents.data = append(pdf.page.contents.data, []byte(instructions)...)
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

	instructions += fmt.Sprintf("%-60s %% draw ellipse part 1\n", fmt.Sprintf("%.2F %.2F m %.2F %.2F %.2F %.2F %.2F %.2F c", (x+rx)*k, (y)*k, (x+rx)*k, (y-ly)*k, (x+lx)*k, (y-ry)*k, x*k, (y-ry)*k))
	instructions += fmt.Sprintf("%-60s %% draw ellipse part 2\n", fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", (x-lx)*k, (y-ry)*k, (x-rx)*k, (y-ly)*k, (x-rx)*k, (y)*k))
	instructions += fmt.Sprintf("%-60s %% draw ellipse part 3\n", fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", (x-rx)*k, (y+ly)*k, (x-lx)*k, (y+ry)*k, x*k, (y+ry)*k))
	instructions += fmt.Sprintf("%-60s %% draw ellipse part 4\n", fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c %s", (x+lx)*k, (y+ry)*k, (x+rx)*k, (y+ly)*k, (x+rx)*k, (y)*k, op))

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

func (pdf *Pdf) ImportPage(sourceFile string, pageno int, box string) int {
	var tplid int

	pdf.fpdi.SetSourceFile(sourceFile)
	pdf.fpdi.SetNextObjectID(pdf.n + 1)

	tplid = pdf.fpdi.ImportPage(pageno, box)

	// write imported objects
	for tplName, objId := range pdf.fpdi.PutFormXobjects() {
		pdf.tplObjIds[tplName] = objId
	}

	// write objects
	objs := pdf.fpdi.GetImportedObjects()
	for i := pdf.n; i < len(objs)+pdf.n; i++ {
		if objs[i] != "" {
			pdf.newObjId()
			//panic(fmt.Sprintf("new object: %d", pdf.n))
			pdf.newObj(pdf.n)
			pdf.outln(objs[i])
		}
	}

	return tplid
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
	pdf.outln(string(pdf.page.contents.data))
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
