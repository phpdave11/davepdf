package davepdf

import "fmt"

func (pdf *Pdf) ImportPage(sourceFile string, pageno int, box string) int {
	var tplid int

	pdf.fpdi.SetSourceFile(sourceFile)
	pdf.fpdi.SetNextObjectID(pdf.n + 1)

	tplid = pdf.fpdi.ImportPage(pageno, box)

	// write imported objects
	for tplName, objId := range pdf.fpdi.PutFormXobjects() {
		pdf.tplObjIds[tplName] = objId
	}

	// write objects
	objs := pdf.fpdi.GetImportedObjects()
	for i := pdf.n; i < len(objs)+pdf.n; i++ {
		if objs[i] != "" {
			pdf.newObjId()
			pdf.newObj(pdf.n)
			pdf.outln(objs[i])
		}
	}

	return tplid
}

func (pdf *Pdf) UseImportedTemplate(tplid int, x, y, w, h float64) {
	// get gofpdi template values
	tplName, scaleX, scaleY, tX, tY := pdf.fpdi.UseTemplate(tplid, x, y, w, h)
	pdf.page.instructions.add(fmt.Sprintf("q 0 J 1 w 0 j 0 G 0 g q %.4F 0 0 %.4F %.4F %.4F cm %s Do Q Q", scaleX, scaleY, tX, tY, tplName), "draw template")
}
