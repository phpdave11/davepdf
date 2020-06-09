package davepdf

import (
	"bytes"
	"encoding/ascii85"
	"fmt"
	"io/ioutil"
	"strconv"
)

type PdfImage struct {
	id       int
	family   string
	contents string
	w        int
	h        int
}

func (pdf *Pdf) newImage() *PdfImage {
	image := &PdfImage{}

	pdf.newObjId()
	image.id = pdf.n

	pdf.images = append(pdf.images, image)

	return image
}

func (pdf *Pdf) NewJPEGImageFromFile(filename string) *PdfImage {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Encode JPEG with ASCII85
	out := bytes.NewBuffer(nil)
	enc := ascii85.NewEncoder(out)
	enc.Write(content)
	enc.Close()
	text := out.String()

	image := pdf.newImage()
	//image.w = 256
	//image.h = 256
	image.contents = text

	return image
}

func (pdf *Pdf) DrawImage(image *PdfImage, x, y float64, w, h int) {
	pdf.page.instructions.add(fmt.Sprintf("q %d 0 0 %d %.5f %.5f cm /Image%d Do Q", w, h, x, pdf.h-float64(h)-y, image.id), "draw image")
}

func (pdf *Pdf) writeImages() {
	for _, image := range pdf.images {
		pdf.newObj(image.id)
		pdf.outln("<<")
		pdf.outln("  /Type /XObject")
		pdf.outln("  /Subtype /Image")
		pdf.outln("  /Name /Image" + strconv.Itoa(image.id))
		pdf.outln("  /Width 256")
		pdf.outln("  /Height 256")
		pdf.outln("  /ColorSpace /DeviceRGB")
		pdf.outln("  /BitsPerComponent 8")
		pdf.outln("  /Length 15684")
		pdf.outln("  /Filter [/ASCII85Decode /DCTDecode]")
		pdf.outln(">>")
		pdf.outln("stream")
		pdf.outln(image.contents)
		pdf.outln("endstream")
		pdf.outln("endobj\n")
	}
}
