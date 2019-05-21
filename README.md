# davepdf

## pdf generator library

### example

```
package main

import "github.com/phpdave11/davepdf"

func main() {
    pdf := davepdf.NewPdf()

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

    pdf.SetFillColorCMYK(0, 81, 81, 45)
    pdf.Rect(10, 200, 250, 50, "F")

    pdf.SetFillColorCMYK(26, 0, 99, 13)
    pdf.Circle(110, 300, 70, "F")

    pdf.Write()
}
```
