package davepdf

import (
	"fmt"
	"math"
)

func (pdf *Pdf) Ellipse(x, y, rx, ry float64, style string) {
	var op string
	var lx, ly, k float64

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

	pdf.page.instructions.add(fmt.Sprintf("%.2F %.2F m %.2F %.2F %.2F %.2F %.2F %.2F c", (x+rx)*k, (pdf.h-y)*k, (x+rx)*k, (pdf.h-(y-ly))*k, (x+lx)*k, (pdf.h-(y-ry))*k, x*k, (pdf.h-(y-ry))*k), "draw ellipse part 1")
	pdf.page.instructions.add(fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", (x-lx)*k, (pdf.h-(y-ry))*k, (x-rx)*k, (pdf.h-(y-ly))*k, (x-rx)*k, (pdf.h-y)*k), "draw ellipse part 2")
	pdf.page.instructions.add(fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", (x-rx)*k, (pdf.h-(y+ly))*k, (x-lx)*k, (pdf.h-(y+ry))*k, x*k, (pdf.h-(y+ry))*k), "draw ellipse part 3")
	pdf.page.instructions.add(fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c %s", (x+lx)*k, (pdf.h-(y+ry))*k, (x+rx)*k, (pdf.h-(y+ly))*k, (x+rx)*k, (pdf.h-y)*k, op), "draw ellipse part 4")
}
