package davepdf

import (
	"fmt"
	"strings"
)

type FunctionType = int

const (
	FunctionType0 FunctionType = 0
	FunctionType2              = 2
	FunctionType3              = 3
	FunctionType4              = 4
)

type PdfFunction struct {
	id            int
	Type          FunctionType
	Domain        []float64
	Size          []int
	Encode        []float64
	BitsPerSample int
	Range         []float64
	Decode        []float64
	Bounds        []float64
	Functions     []*PdfFunction
	C0            []float64
	C1            []float64
	N             float64
}

func (pdf *Pdf) NewFunction() *PdfFunction {
	function := &PdfFunction{}

	pdf.newObjId()
	function.id = pdf.n

	pdf.functions = append(pdf.functions, function)

	return function
}

func (pdf *Pdf) writeFunctions() {
	for _, function := range pdf.functions {
		childFuncs := ""
		for _, childFunc := range function.Functions {
			childFuncs += fmt.Sprintf("%d 0 R ", childFunc.id)
		}
		// Trim leading space
		childFuncs = strings.TrimSpace(childFuncs)

		pdf.newObj(function.id)
		pdf.outln("<<")
		pdf.outln(fmt.Sprintf("  /FunctionType %d", function.Type))
		pdf.outln(fmt.Sprintf("  /Domain [%f %f]", function.Domain[0], function.Domain[1]))

		if function.Type == FunctionType3 {
			pdf.outln(fmt.Sprintf("  /Functions [%s]", childFuncs))
			pdf.outln(fmt.Sprintf("  /Bounds [%f]", function.Bounds[0]))
			pdf.outln(fmt.Sprintf("  /Encode [%f %f %f %f]", function.Encode[0], function.Encode[1], function.Encode[2], function.Encode[3]))
		}

		if function.Type == FunctionType2 {
			pdf.outln(fmt.Sprintf("  /C0 [%f %f %f %f]", function.C0[0], function.C0[1], function.C0[2], function.C0[3]))
			pdf.outln(fmt.Sprintf("  /C1 [%f %f %f %f]", function.C1[0], function.C1[1], function.C1[2], function.C1[3]))
			pdf.outln(fmt.Sprintf("  /N %f", function.N))
		}

		pdf.outln(">>")
		pdf.outln("endobj\n")
	}
}
