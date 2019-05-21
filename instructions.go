package davepdf

import "fmt"

type PdfInstruction struct {
	code     string
	comments string
}

type PdfInstructions struct {
	instructions []*PdfInstruction
}

func (ins *PdfInstructions) add(code string, comments string) {
	in := &PdfInstruction{code: code, comments: comments}
	ins.instructions = append(ins.instructions, in)
}

func (ins *PdfInstructions) String() string {
	result := ""

	maxLen := 0
	for i := 0; i < len(ins.instructions); i++ {
		// Get instruction
		instruction := ins.instructions[i]

		// Set max length
		if len(instruction.code) > maxLen {
			maxLen = len(instruction.code)
		}
	}

	for i := 0; i < len(ins.instructions); i++ {
		// Get instruction
		instruction := ins.instructions[i]

		// Append result
		result += fmt.Sprintf("%-"+fmt.Sprintf("%d", maxLen+2)+"s %% %s\n", instruction.code, instruction.comments)
	}

	return result
}

func (pdf *Pdf) newInstructions() *PdfInstructions {
	result := &PdfInstructions{}
	result.instructions = make([]*PdfInstruction, 0)
	return result
}
