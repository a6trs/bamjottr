package grass

import "testing"

type navDispTest struct {
	curpage, pagecnt int
	expected         []int
}

// Automated tests for navigationDisplay [grass/projects.go]
func TestNavDisp(t *testing.T) {
	tests := []navDispTest{
		{ // #1
			1, 1,
			[]int{1},
		},
		{ // #2
			2, 3,
			[]int{1, 2, 3},
		},
		{ // #3
			1, 6,
			[]int{1, 2, 3, -1, 6},
		},
		{ // #4
			2, 6,
			[]int{1, 2, 3, 4, 5, 6},
		},
		{ // #5
			2, 10,
			[]int{1, 2, 3, 4, -1, 10},
		},
		{ // #6
			6, 10,
			[]int{1, -1, 4, 5, 6, 7, 8, 9, 10},
		},
		{ // #7
			9, 10,
			[]int{1, -1, 7, 8, 9, 10},
		},
		{ // #8
			9, 15,
			[]int{1, -1, 7, 8, 9, 10, 11, -1, 15},
		},
	}
	for i, test := range tests {
		doTestCaseNavDisp(t, i, test)
	}
}

func doTestCaseNavDisp(t *testing.T, tcnum int, test navDispTest) {
	result := navigationDisplay(test.curpage, test.pagecnt)
	if len(result) != len(test.expected) {
		t.Errorf("(%d) Length different: expected %d, got %d", tcnum+1, len(test.expected), len(result))
	}
	for i, r := range test.expected {
		if result[i] != r {
			t.Errorf("(%d) Result not match\n\tExpected %v\n\tGot      %v", tcnum+1, test.expected, result)
			break
		}
	}
}
