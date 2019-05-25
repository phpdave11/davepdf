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
    //pdf.AddPage()

    pdf.SetFontFamily("Times-Roman")
    pdf.SetFontSize(28)
    pdf.SetXY(10, 600)
    pdf.SetXY(0, 0)
    pdf.Text("Hello World!")

    pdf.SetFillColorCMYK(0, 81, 81, 45)
    pdf.Rect(140, 140, 50, 50, "F")

    pdf.SetFillColorCMYK(26, 0, 99, 13)
    pdf.Circle(70, 140, 70, "F")

    pdf.SetLineWidth(5)
    pdf.SetLineCapStyle(1)
    pdf.SetLineJoinStyle(0)
    pdf.SetLineDashPattern([]uint{10, 10}, 0)
    pdf.SetStrokeColorCMYK(20, 20, 100, 70)
    pdf.SetFillColorCMYK(90, 20, 80, 0)
    pdf.NewPath(400.0, 400.0)
    pdf.AppendCurve(300.0, 450.0, 360.0, 380.0, 400.0, 420.0)
    pdf.StrokeAndFillPath()

    pdf.Write()
}

/*
%1.42 w                                         % set line width
%0 J                                            % set line cap style
%0 j                                            % set line join style
%[10.00 10.00] 0.00 d                           % set line dash pattern
%0.000 1.000 0.000 RG                           % set rgb color for stroking operation
%14.17 728.50 m                                 % begin a new submath by moving current point to (x, y)
%85.04 685.98 198.43 714.33 170.08 629.29 c     % append a cubic bezier curve to the current path
%S                                              % stroke the path

226.77 728.50 m                                % begin a new submath by moving current point to (x, y)
198.43 629.29 425.20 714.33 283.46 629.29 c    % append a cubic bezier curve to the current path
f                                              % fill the path with nonzero winding number rule

0.784 0.863 0.784 rg
1.42 w
0 J
0 j
[10.00 10.00] 0.00 d
0.000 1.000 0.000 RG
396.85 728.50 m
425.20 685.98 510.24 714.33 566.93 629.29 c
B
*/

```
