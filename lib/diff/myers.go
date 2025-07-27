package diff

import (
	"geo-git/lib"
	"geo-git/lib/utils"
	"slices"
)

type MyersDiff struct {
	Diff *lib.Diff
}

func NewMyersDiff(Diff *lib.Diff) *MyersDiff {
	temp := &MyersDiff{Diff: Diff}
	return temp
}

func (md *MyersDiff) DoDiff() {
	for _, backTrack := range md.backTrack() {
		var beforeLine string
		if len(md.Diff.LinesBefore) <= backTrack.prev_x {
			beforeLine = ""
		} else {
			beforeLine = md.Diff.LinesBefore[backTrack.prev_x]
		}
		var afterLine string
		if len(md.Diff.LinesAfter) <= backTrack.prev_y {
			afterLine = ""
		} else {
			afterLine = md.Diff.LinesAfter[backTrack.prev_y]
		}

		if backTrack.x == backTrack.prev_x {
			md.Diff.Edits = append(md.Diff.Edits, lib.NewEdit(lib.Ins, afterLine))
		} else if backTrack.y == backTrack.prev_y {
			md.Diff.Edits = append(md.Diff.Edits, lib.NewEdit(lib.Del, beforeLine))
		} else {
			md.Diff.Edits = append(md.Diff.Edits, lib.NewEdit(lib.Eql, beforeLine))
		}
	}
	slices.Reverse(md.Diff.Edits)
}

type backTrack struct {
	prev_x, prev_y, x, y int
}

func (md *MyersDiff) backTrack() []backTrack {

	backTracks := []backTrack{}
	x := len(md.Diff.LinesBefore)
	y := len(md.Diff.LinesAfter)

	shortest_edit := md.shortestEdit()
	if shortest_edit == nil {
		return backTracks
	}
	for d := len(shortest_edit) - 1; d >= 0; d-- {
		v := shortest_edit[d]
		k := x - y
		var prev_k int
		if k == -d || (k != d && utils.GetElement(v, k-1) < utils.GetElement(v, k+1)) {
			prev_k = k + 1
		} else {
			prev_k = k - 1
		}

		prev_x := utils.GetElement(v, prev_k)
		prev_y := prev_x - prev_k

		for {
			if !(x > prev_x && y > prev_y) {
				break
			}
			bk := backTrack{prev_x: x - 1, prev_y: y - 1, x: x, y: y}
			backTracks = append(backTracks, bk)
			x = x - 1
			y = y - 1
		}
		if d > 0 {
			bk := backTrack{prev_x: prev_x, prev_y: prev_y, x: x, y: y}
			backTracks = append(backTracks, bk)
		}
		x, y = prev_x, prev_y
	}
	return backTracks
}

func (md *MyersDiff) shortestEdit() [][]int {
	n := len(md.Diff.LinesBefore)
	m := len(md.Diff.LinesAfter)
	max := n + m

	v := make([]int, max*2+1)
	v[1] = 0
	trace := [][]int{}
	for d := 0; d <= max; d++ {
		copy_v := make([]int, len(v))
		copy(copy_v, v)
		trace = append(trace, copy_v)
		for k := -1 * d; k <= d; k = k + 2 {
			var x int
			if k == -d || (k != d && utils.GetElement(v, k-1) < utils.GetElement(v, k+1)) {
				x = utils.GetElement(v, k+1)
			} else {
				x = utils.GetElement(v, k-1) + 1
			}
			y := x - k

			for {

				if !(x < n && y < m && md.Diff.LinesBefore[x] == md.Diff.LinesAfter[y]) {
					break
				}
				x = x + 1
				y = y + 1
			}
			utils.SaveElement(v, k, x)
			if x >= n && y >= m {
				return trace
			}
		}
	}
	return trace
}
