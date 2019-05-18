# davepdf

## pdf generator library

### example

```
package main

import "github.com/phpdave11/davepdf"

func main() {
    pdf := davepdf.NewPdf()

    // add page
    pdf.AddPage()

    tpl1 := pdf.ImportPage("/tmp/PDFPL110.pdf", 1, "/BleedBox")
    pdf.UseImportedTemplate(tpl1, 300, -450, 250, 0)

    tpl2 := pdf.ImportPage("/tmp/PDFPL114.pdf", 1, "/BleedBox")
    pdf.UseImportedTemplate(tpl2, 300, -250, 250, 0)

    pdf.SetFontFamily("Times-Roman")
    pdf.SetFontSize(28)
    pdf.SetXY(10, 600)
    pdf.SetXY(0, 0)
    pdf.Text("Hello Courtney!  I love you!")

    //pdf.SetFillColor(&CMYKColor{C: 48, M: 32, Y: 0, K: 0})
    pdf.SetFillColor(&davepdf.CMYKColor{C: 0, M: 81, Y: 81, K: 45})
    pdf.Rect(10, 200, 250, 50, "F")

    pdf.SetFillColor(&davepdf.CMYKColor{C: 26, M: 0, Y: 99, K: 13})
    //pdf.Ellipse(100, 50, 30, 20, "D")
    pdf.Circle(110, 300, 70, "F")

    pdf.Write()
}
```