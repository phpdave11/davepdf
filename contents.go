package davepdf

import "fmt"

type PdfContents struct {
	id   int
	data []byte
}

func (pdf *Pdf) newContents() *PdfContents {
	contents := &PdfContents{}

	pdf.newObjId()
	contents.id = pdf.n

	return contents
}

func (pdf *Pdf) writeContents() {
	for _, page := range pdf.pageTree.pages {
		inst := page.instructions.String()

		inst += `
316.789  140.311  m % Move to start of leaf
303.222  146.388  282.966  136.518  279.122  121.983  c % Curved segment
277.322  120.182  l % Straight line
285.125  122.688  291.441  121.716  298.156  119.386  c % Curved segment
336.448  119.386  l % Straight line
331.072  128.643  323.346  137.376  316.789  140.311  c % Curved segment
W n % Set clipping path
q % Save graphics state
27.7843  0.0000  0.0000 −27.7843  310.2461  121.1521  cm % Set matrix
/Sh1  sh % Paint shading
Q % Restore graphics state
`

		pdf.newObj(page.contents.id)
		pdf.outln("<<")
		pdf.outln(fmt.Sprintf("  /Length %d", len(inst)))
		pdf.outln(">>")
		pdf.outln("stream")
		pdf.outln(inst)
		pdf.outln("endstream")
		pdf.outln("endobj\n")
	}
}
