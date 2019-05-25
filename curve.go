package davepdf

import (
	"fmt"
)

func (pdf *Pdf) NewPath(x, y float64) {
	pdf.page.instructions.add(fmt.Sprintf("%.2F %.2F m", x*pdf.k, (pdf.h-y)*pdf.k), "begin a new subpath at (x, y)")
}

func (pdf *Pdf) StrokePath() {
	pdf.page.instructions.add("S", "stroke the path")
}

func (pdf *Pdf) FillPath() {
	pdf.page.instructions.add("f", "fill the path")
}

func (pdf *Pdf) StrokeAndFillPath() {
	pdf.page.instructions.add("B", "stroke AND fill the path")
}

func (pdf *Pdf) AppendCurve(x1, y1, x2, y2, x3, y3 float64) {
	pdf.page.instructions.add(fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", x1*pdf.k, (pdf.h-y1)*pdf.k, x2*pdf.k, (pdf.h-y2)*pdf.k, x3*pdf.k, (pdf.h-y3)*pdf.k), "draw a Bézier curve from last draw point")
}
