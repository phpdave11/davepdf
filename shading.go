package davepdf

import (
	"fmt"
)

type ShadingType = int

const (
	ShadingType3 ShadingType = 3
)

type PdfShading struct {
	id       int
	Type     ShadingType
	Coords   []float64
	Function *PdfFunction
	Extend   []bool
}

func (pdf *Pdf) NewShading() *PdfShading {
	shading := &PdfShading{}

	pdf.newObjId()
	shading.id = pdf.n

	pdf.shadings = append(pdf.shadings, shading)

	return shading
}

func (pdf *Pdf) writeShadings() {
	for _, shading := range pdf.shadings {
		pdf.newObj(shading.id)
		pdf.outln("<<")
		pdf.outln("  /ColorSpace /DeviceCMYK")
		pdf.outln(fmt.Sprintf("  /ShadingType %d", shading.Type))
		pdf.outln(fmt.Sprintf("  /Coords [%f %f %f %f %f %f]", shading.Coords[0], shading.Coords[1], shading.Coords[2], shading.Coords[3], shading.Coords[4], shading.Coords[5]))
		pdf.outln(fmt.Sprintf("  /Function %d 0 R", shading.Function.id))
		pdf.outln(fmt.Sprintf("  /Extend %t", shading.Extend))

		pdf.outln(">>")
		pdf.outln("endobj\n")
	}
}
