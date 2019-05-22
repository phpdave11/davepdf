package davepdf

import "fmt"

func (pdf *Pdf) Text(text string) {
	pdf.page.instructions.add("BT", "begin text")
	pdf.page.instructions.add(fmt.Sprintf("  %s %d Tf", pdf.fontFamily, pdf.fontSize), "set font family and font size")
	pdf.page.instructions.add(fmt.Sprintf("  %f %f Td", pdf.x, pdf.y), "set position to draw text")
	pdf.page.instructions.add(fmt.Sprintf("  (%s)Tj", text), "write text")
	pdf.page.instructions.add("ET", "end text")
}
