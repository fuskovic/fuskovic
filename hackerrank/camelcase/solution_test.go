package camelcase

import "testing"

func TestSolution(t *testing.T) {
	for _, test := range []struct {
		word string
		want int
	}{
		{
			word: "oneTwoThree",
			want: 3,
		},
		{
			word: "saveChangesInTheEditor",
			want: 5,
		},
	} {
		t.Run(test.word, func(t *testing.T) {
			if got := solution(test.word); got != test.want {
				t.Errorf("got %d but wanted %d\n", got, test.want)
			}
		})
	}
}
