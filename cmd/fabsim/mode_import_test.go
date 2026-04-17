package main

import (
	"bufio"
	"strings"
	"testing"
)

// TestReadUntilFabraryFooter pins the paste-terminator behaviour: input ends at the fabrary footer
// line, content after the footer is ignored, and EOF without a footer is still accepted.
func TestReadUntilFabraryFooter(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "stops at footer, drops trailing junk",
			in: "Hero: Viserai\n" +
				"Deck cards\n" +
				"2x Rune Flash (red)\n" +
				"Made with ❤️ at the FaBrary\n" +
				"See the full deck @ https://fabrary.net/decks/ABC\n" +
				"GARBAGE AFTER FOOTER\n",
			want: "Hero: Viserai\n" +
				"Deck cards\n" +
				"2x Rune Flash (red)\n" +
				"Made with ❤️ at the FaBrary\n" +
				"See the full deck @ https://fabrary.net/decks/ABC\n",
		},
		{
			name: "no footer — EOF ends input",
			in:   "Hero: Viserai\n2x Rune Flash (red)\n",
			want: "Hero: Viserai\n2x Rune Flash (red)\n",
		},
		{
			name: "footer on final line without trailing newline",
			in:   "Hero: Viserai\nSee the full deck @ https://fabrary.net/decks/XYZ",
			want: "Hero: Viserai\nSee the full deck @ https://fabrary.net/decks/XYZ",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := readUntilFabraryFooter(bufio.NewReader(strings.NewReader(c.in)))
			if err != nil {
				t.Fatalf("readUntilFabraryFooter: %v", err)
			}
			if string(got) != c.want {
				t.Errorf("got:\n%q\nwant:\n%q", got, c.want)
			}
		})
	}
}
