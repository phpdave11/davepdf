package davepdf

import (
	"fmt"
)

func (pdf *Pdf) Curve(x0, y0, x1, y1, x2, y2, x3, y3 float64, style string) {
	var op string

	if style == "F" {
		op = "f"
	} else if style == "FD" || style == "DF" {
		op = "B"
	} else {
		op = "S"
	}

	if op == "123" {
		fmt.Printf("")
	}

	pdf.point(x0, y0)
	pdf.curve(x1, y1, x2, y2, x3, y3)
}

func (pdf *Pdf) point (x, y float64) {
    // Sets a draw point
    // Parameters:
    // - x, y: Point
	pdf.page.instructions.add(fmt.Sprintf("%.2F %.2F m", x * pdf.k, (pdf.h - y) * pdf.k), "set draw point")
}

func (pdf *Pdf) curve(x1, y1, x2, y2, x3, y3 float64) {
    // Draws a Bézier curve from last draw point
    // Parameters:
    // - x1, y1: Control point 1
    // - x2, y2: Control point 2
    // - x3, y3: End point
	pdf.page.instructions.add(fmt.Sprintf("%.2F %.2F %.2F %.2F %.2F %.2F c", x1 * pdf.k, (pdf.h - y1) * pdf.k, x2 * pdf.k, (pdf.h - y2) * pdf.k, x3 * pdf.k, (pdf.h - y3) * pdf.k), "draw a Bézier curve from last draw point")
}
