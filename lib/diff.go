package lib

import "strings"

type Diff struct {
	DocumentBefore string
	LinesBefore    []string
	DocumentAfter  string
	LinesAfter     []string
	Edits          []*Edit
}

type Edit struct {
	editType Symbol
	text     string
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
	diff.LinesBefore = strings.Split(before, "\n")
	diff.LinesAfter = strings.Split(after, "\n")
	diff.Edits = []*Edit{}
	return diff
}

func NewEdit(editType Symbol, text string) *Edit {
	return &Edit{editType: editType, text: text}
}

func (ed *Edit) ToString() string {
	stringSymbol, ok := SymbolMaps[ed.editType]
	if ok {
		return stringSymbol + " " + ed.text
	}
	return ""
}
