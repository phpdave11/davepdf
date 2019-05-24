# davepdf

## pdf generator library

### example

```
package main

import "github.com/phpdave11/davepdf"

func main() {
    pdf := davepdf.NewPdf()

	// Page 1
    pdf.AddPage()

    tpl1 := pdf.ImportPage("/tmp/PDFPL110.pdf", 1, "/BleedBox")
    pdf.UseImportedTemplate(tpl1, 362, 0, 250, 0)

	// Page 2
    pdf.AddPage()

    pdf.SetFontFamily("Times-Roman")
    pdf.SetFontSize(28)
    pdf.SetXY(10, 600)
    pdf.SetXY(0, 0)
    pdf.Text("Hello World!")

    pdf.SetFillColorCMYK(0, 81, 81, 45)
    pdf.Rect(140, 140, 50, 50, "F")

    pdf.SetFillColorCMYK(26, 0, 99, 13)
    pdf.Circle(70, 140, 70, "F")

	// Curve (does not work yet)
	pdf.Curve(5, 40, 30, 55, 70, 45, 60, 75, "")
	pdf.Curve(80, 40, 70, 75, 150, 45, 100, 75, "")
	pdf.Curve(140, 40, 150, 55, 180, 45, 200, 75, "")

    pdf.Write()
}
```
