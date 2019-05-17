package davepdf

import (
	"bytes"
	"fmt"
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
	page       *PdfPage
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

func (pdf *Pdf) SetFontFamily(fontFamily string) {
	pdf.fontFamily = "/FONT1"
}

func (pdf *Pdf) SetFontSize(fontSize int) {
	pdf.fontSize = fontSize
}

func (pdf *Pdf) SetXY(x, y float64) {
	pdf.x = x
	pdf.y = y
}

func (pdf *Pdf) Text(text string) {
	var instructions string

	instructions += fmt.Sprintf("  %-30s %% begin text\n", "BT")
	instructions += fmt.Sprintf("    %-28s %% set font family and font size\n", fmt.Sprintf("%s %d Tf", pdf.fontFamily, pdf.fontSize))
	instructions += fmt.Sprintf("    %-28s %% set position to draw text\n", fmt.Sprintf("%f %f Td", pdf.x, pdf.y))
	instructions += fmt.Sprintf("    %-28s %% write text\n", fmt.Sprintf("(%s)Tj", text))
	instructions += fmt.Sprintf("  %-30s %% end text", "ET")

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
	pdf.outln("  /BaseFont /Helvetica")
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

	pdf.SetFontFamily("Helvetica")
	pdf.SetFontSize(18)
	pdf.SetXY(10, 600)
	pdf.SetXY(0, 0)
	pdf.Text("Hello World!")

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
