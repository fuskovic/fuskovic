package caesarcipher

import "testing"

func TestSolution(t *testing.T) {
	for _, test := range []struct {
		word, want   string
		spacesToMove int
	}{
		{
			word:         "abcdefghijklmnopqrstuvwxyz",
			spacesToMove: 3,
			want:         "defghijklmnopqrstuvwxyzabc",
		},
		{
			word:         "There's-a-starman-waiting-in-the-sky",
			spacesToMove: 3,
			want:         "Wkhuh'v-d-vwdupdq-zdlwlqj-lq-wkh-vnb",
		},
		{
			word:         "middle-Outz",
			spacesToMove: 2,
			want:         "okffng-Qwvb",
		},
	} {
		t.Run(test.word, func(t *testing.T) {
			if got := solution(test.word, test.spacesToMove); got != test.want {
				t.Errorf("got %s but wanted %s\n", got, test.want)
			}
		})
	}
}
