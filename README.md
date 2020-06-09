# davepdf

## PDF generator library

### Features

- Creates PDF 1.4 documents
- Less than 1000 lines of code (currently: 871)
- Less than 100 lines of code per file
- ASCII text support
- Support for 14 fonts: Times-Roman, Times-Bold, Times-Italic, Times-BoldItalic, Courier, Courier-Bold, Courier-Oblique, Courier-BoldOblique, Helvetica, Helvetica-Bold, Helvetica-Oblique, Helvetica-BoldOblique, Symbol, ZapfDingbats
- JPEG image support
- PDF output contains only ASCII characters and can be opened in Notepad (except when importing existing PDFs)
- PDF output includes comments about each PDF drawing instruction
- Basic shape support: circle, ellipse, rectangle, line, bezier curves
- Import existing PDF page and print on new PDF page

### Limitations

- Unsupported - use at your own risk
- No support for external fonts
- No UTF-8 support

### Example

```
package main

import "github.com/phpdave11/davepdf"

func main() {
    pdf := davepdf.NewPdf()

    pdf.AddPage()

	pdf.SetFontFamily("Courier-Bold")
	pdf.SetFontSize(24)
	pdf.SetXY(180, 10)
	pdf.Text("phpdave11/davepdf")

	fonts := []string{
		"Times-Roman", 
		"Times-Bold", 
		"Times-Italic", 
		"Times-BoldItalic", 
		"Courier", 
		"Courier-Bold", 
		"Courier-Oblique", 
		"Courier-BoldOblique", 
		"Helvetica", 
		"Helvetica-Bold", 
		"Helvetica-Oblique", 
		"Helvetica-BoldOblique", 
		"Symbol", 
		"ZapfDingbats"}

	for i := 0; i < len(fonts); i++ {
		pdf.SetFontFamily(fonts[i])
		pdf.SetFontSize(20)
		pdf.SetXY(10, 80 + float64(i * 25))
		pdf.Text("The quick brown fox jumps over the lazy dog.")
	}

	image := pdf.NewJPEGImageFromFile("/Users/dave-barnes/Desktop/go-img/go.jpg")
	pdf.DrawImage(image, 178, 460, 256, 256)

    pdf.Write()
}
```

Output (rendered PDF):


Output (PDF drawing instructions):

```
BT                                                   % begin text
  /FONT6 24 Tf                                       % set font family and font size
  180.000000 758.000000 Td                           % set position to draw text
  (phpdave11/davepdf)Tj                              % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT7 20 Tf                                       % set font family and font size
  10.000000 692.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT8 20 Tf                                       % set font family and font size
  10.000000 667.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT9 20 Tf                                       % set font family and font size
  10.000000 642.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT10 20 Tf                                      % set font family and font size
  10.000000 617.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT11 20 Tf                                      % set font family and font size
  10.000000 592.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT6 20 Tf                                       % set font family and font size
  10.000000 567.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT12 20 Tf                                      % set font family and font size
  10.000000 542.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT13 20 Tf                                      % set font family and font size
  10.000000 517.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT14 20 Tf                                      % set font family and font size
  10.000000 492.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT15 20 Tf                                      % set font family and font size
  10.000000 467.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT16 20 Tf                                      % set font family and font size
  10.000000 442.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT17 20 Tf                                      % set font family and font size
  10.000000 417.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT18 20 Tf                                      % set font family and font size
  10.000000 392.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
BT                                                   % begin text
  /FONT19 20 Tf                                      % set font family and font size
  10.000000 367.000000 Td                            % set position to draw text
  (The quick brown fox jumps over the lazy dog.)Tj   % write text
ET                                                   % end text
q 256 0 0 256 178.00000 76.00000 cm /Image20 Do Q    % draw image
```
