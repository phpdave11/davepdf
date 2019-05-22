package davepdf

import "fmt"

func (pdf *Pdf) Rect(x, y, w, h float64, mode string) {
	pdf.page.instructions.add(fmt.Sprintf("%.5f %.5f %.5f %.5f re %s", x, y, w, h, mode), "draw rectangle")
}
