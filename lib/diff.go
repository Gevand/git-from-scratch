package lib

import (
	"fmt"
	"strings"
)

type Diff struct {
	DocumentBefore string
	LinesBefore    []Line
	DocumentAfter  string
	LinesAfter     []Line
	Edits          []Edit
}

type Edit struct {
	editType   Symbol
	beforeLine Line
	afterLine  Line
}

type Line struct {
	number int
	Text   string
}

const HUNK_CONTEXT = 3

type Hunk struct {
	BeforeStart int
	AfterStart  int
	Edits       []Edit
}

func GetDummyLine() Line {
	return Line{Text: "Dummy", number: -1}
}

type Symbol int

const (
	Eql Symbol = iota
	Ins
	Del
)

var SymbolMaps = map[Symbol]string{
	Eql: " ",
	Ins: "+",
	Del: "-",
}

func NewDiff(before string, after string) *Diff {
	diff := &Diff{DocumentBefore: before, DocumentAfter: after}
	diff.LinesBefore = lines(before)
	diff.LinesAfter = lines(after)
	diff.Edits = []Edit{}
	return diff
}

func NewHunk(before, after int, edits []Edit) *Hunk {
	return &Hunk{BeforeStart: before, AfterStart: after, Edits: edits}
}

func NewEdit(editType Symbol, beforeLine Line, afterLine Line) *Edit {
	return &Edit{editType: editType, beforeLine: beforeLine, afterLine: afterLine}
}

func (hunk *Hunk) Filter() []Hunk {
	hunks := []Hunk{}
	offset := 0
	for {
		for {
			offset += 1
			if offset > len(hunk.Edits) || !(hunk.Edits[offset].editType == Eql) {
				break
			}
		}
		if offset >= len(hunk.Edits) {
			return hunks
		}
		offset = offset - (HUNK_CONTEXT + 1)
		if offset < 0 {
			hunk.BeforeStart = 0
			hunk.AfterStart = 0
		} else {
			hunk.BeforeStart = hunk.Edits[offset].beforeLine.number
			hunk.AfterStart = hunk.Edits[offset].afterLine.number
		}
		new_hunk := *NewHunk(hunk.AfterStart, hunk.AfterStart, []Edit{})
		offset = buildHunk(&new_hunk, hunk.Edits, offset)
		hunks = append(hunks, new_hunk)
	}

}
func buildHunk(hunk *Hunk, edits []Edit, offset int) int {
	counter := -1
	for {
		if counter == 0 {
			break
		}
		if offset > 0 && counter > 0 {
			hunk.Edits = append(hunk.Edits, edits[offset])
		}
		offset = offset + 1
		if offset >= len(edits) {
			break
		}
		if offset+HUNK_CONTEXT < len(edits) {
			t := edits[offset+HUNK_CONTEXT].editType
			if t == Ins || t == Del {
				counter = 2*HUNK_CONTEXT + 1
			} else {
				counter = counter - 1
			}
		}
	}
	return offset
}

func (hunk *Hunk) GenerateHeader() string {

	before_offset_start := -1
	after_offset_start := -1
	before_offset_line_size := 0
	after_offset_line_size := 0
	dummy := GetDummyLine()
	for _, edit := range hunk.Edits {
		if edit.beforeLine.Text != dummy.Text {
			if before_offset_start != -1 {
				before_offset_start = edit.beforeLine.number
			}
			before_offset_line_size = before_offset_line_size + 1
		}
		if edit.afterLine.Text != dummy.Text {
			if after_offset_start != -1 {
				after_offset_start = edit.afterLine.number
			}
			after_offset_line_size = after_offset_line_size + 1
		}
	}
	before_offset_start = max(before_offset_start, 0)
	after_offset_start = max(after_offset_start, 0)

	return fmt.Sprintf("@@ -#%d %d +#%d %d @@", before_offset_start, before_offset_line_size, after_offset_start, after_offset_line_size)
}

func (ed *Edit) ToString() string {
	stringSymbol, ok := SymbolMaps[ed.editType]
	line := pickLine(ed.beforeLine, ed.afterLine)
	if ok {
		return stringSymbol + " " + line.Text
	}
	return ""
}

func pickLine(beforeLine, afterLine Line) Line {
	dummy := GetDummyLine()
	if beforeLine.number == dummy.number && beforeLine.Text == dummy.Text {
		return afterLine
	}
	return beforeLine
}

func lines(document string) []Line {
	lines := []Line{}
	line_array := strings.Split(document, "\n")
	for index, line := range line_array {
		temp_line := Line{number: index + 1, Text: line}
		lines = append(lines, temp_line)
	}
	return lines
}
