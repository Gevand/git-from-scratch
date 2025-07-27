package unit_test

import (
	"geo-git/lib"
	"geo-git/lib/diff"
	"strings"
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
	expected_value := "- A - B C + B A B - B A + C"
	expected_value = strings.ReplaceAll(expected_value, "\r", "")
	expected_value = strings.ReplaceAll(expected_value, "\n", "")
	expected_value = strings.ReplaceAll(expected_value, " ", "")
	result = strings.ReplaceAll(result, "\r", "")
	result = strings.ReplaceAll(result, "\n", "")
	result = strings.ReplaceAll(result, " ", "")

	if result != expected_value {
		t.Errorf("\nExpcted: %v - got %v\n", expected_value, result)
	}
}
