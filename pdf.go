package davepdf

import (
	"bytes"
	"fmt"
	"github.com/phpdave11/gofpdi"
)

type Pdf struct {
	n               int
	w               *bytes.Buffer
	catalog         *PdfCatalog
	resources       *PdfResources
	fonts           []*PdfFont
	images          []*PdfImage
	pageTree        *PdfPageTree
	functions       []*PdfFunction
	shadings        []*PdfShading
	xref            *PdfXrefTable
	fpdi            *gofpdi.Importer
	tplObjIds       map[string]int
	offsets         map[int]int
	objects         []*PdfObject
	fontFamily      string
	fontSize        int
	x               float64
	y               float64
	k               float64
	h               float64
	page            *PdfPage
	fillColor       *Color
	strokeColor     *Color
	lineWidth       float64
	lineCapStyle    *LineCap
	lineJoinStyle   *LineJoin
	lineDashPattern *LineDash
}

func NewPdf() *Pdf {
	pdf := &Pdf{}
	pdf.offsets = make(map[int]int, 0)
	pdf.w = new(bytes.Buffer)
	pdf.catalog = pdf.newCatalog()
	pdf.pageTree = pdf.newPageTree()
	pdf.resources = pdf.newResources()
	pdf.catalog.pageTree = pdf.pageTree
	pdf.fpdi = gofpdi.NewImporter()
	pdf.tplObjIds = make(map[string]int, 0)
	pdf.functions = make([]*PdfFunction, 0)
	pdf.shadings = make([]*PdfShading, 0)
	pdf.fonts = make([]*PdfFont, 0)
	pdf.images = make([]*PdfImage, 0)
	pdf.xref = pdf.newXrefTable()

	pdf.k = 1.0
	pdf.h = 792.0

	pdf.outln("%PDF-1.4\n")

	return pdf
}

func (pdf *Pdf) SetXY(x, y float64) {
	pdf.x = x
	pdf.y = y
}

func (pdf *Pdf) outln(s string) {
	pdf.w.WriteString(s)
	pdf.w.WriteString("\n")
}

func (pdf *Pdf) out(s string) {
	pdf.w.WriteString(s)
}

func (pdf *Pdf) outbyte(b byte) {
	pdf.w.WriteByte(b)
}

func (pdf *Pdf) Write() {
	pdf.writeCatalog()
	pdf.writePageTree()
	pdf.writeResources()
	pdf.writeFonts()
	pdf.writeImages()
	pdf.writeFunctions()
	pdf.writeShadings()
	pdf.writePage()
	pdf.writeContents()
	pdf.writeXref()
	pdf.writeTrailer()

	// output PDF
	fmt.Print(string(pdf.w.Bytes()))
}
