package unit_test

import (
	"geo-git/lib"
	"geo-git/lib/diff"
	"testing"
)

func TestDiff_Simple(t *testing.T) {
	before := "A\nB\nC\nA\nB\nB\nA"
	after := "C\nB\nA\nB\nA\nC"
	myersDiff := diff.NewMyersDiff(lib.NewDiff(before, after))
	myersDiff.DoDiff()
	result := ""
	for _, edit := range myersDiff.Diff.Edits {
		result += edit.ToString() + "\n"
	}
	t.Errorf("\n%v\n", result)
}
