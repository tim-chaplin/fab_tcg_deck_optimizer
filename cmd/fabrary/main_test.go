package main

import (
	"bufio"
	"strings"
	"testing"
)

// TestValidateDeckName pins the rejection rules for names that become file paths under mydecks/.
// The concern is path traversal ("../") and Windows-reserved characters that would crash WriteFile.
func TestValidateDeckName(t *testing.T) {
	ok := []string{"viserai-v2", "my best deck", "deck_42", "Viserai"}
	bad := []string{"", ".", "..", "../escape", "path/with/slash", "back\\slash", "has*star", "has?mark"}

	for _, name := range ok {
		if err := validateDeckName(name); err != nil {
			t.Errorf("validateDeckName(%q) = %v, want nil", name, err)
		}
	}
	for _, name := range bad {
		if err := validateDeckName(name); err == nil {
			t.Errorf("validateDeckName(%q) = nil, want error", name)
		}
	}
}

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
