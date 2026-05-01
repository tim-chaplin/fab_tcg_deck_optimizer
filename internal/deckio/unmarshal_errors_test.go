package deckio

import (
	"strings"
	"testing"
)

// Tests that Unmarshal returns descriptive errors when the JSON references unknown heroes,
// weapons, or cards, or when the input is malformed JSON.
func TestUnmarshal_RejectsUnknownNamesAndBadJSON(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantError string
	}{
		{
			name:      "MalformedJSON",
			input:     "{not json",
			wantError: "invalid",
		},
		{
			name:      "UnknownHero",
			input:     `{"hero":"NoSuchHero","weapons":[],"cards":[]}`,
			wantError: `unknown hero "NoSuchHero"`,
		},
		{
			name:      "UnknownWeapon",
			input:     `{"hero":"Viserai","weapons":["Phantom Sword"],"cards":[]}`,
			wantError: `unknown weapon "Phantom Sword"`,
		},
		{
			name:      "UnknownCard",
			input:     `{"hero":"Viserai","weapons":[],"cards":["Made-Up Card"]}`,
			wantError: `unknown card "Made-Up Card"`,
		},
		{
			name:      "UnknownPerCardMarginal",
			input:     `{"hero":"Viserai","weapons":[],"cards":[],"stats":{"per_card_marginal":[{"card":"Made-Up Card"}]}}`,
			wantError: `unknown card "Made-Up Card" in per_card_marginal stats`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Unmarshal([]byte(tc.input))
			if err == nil {
				t.Fatalf("Unmarshal(%q): err = nil, want error containing %q", tc.input, tc.wantError)
			}
			if !strings.Contains(err.Error(), tc.wantError) {
				t.Errorf("Unmarshal(%q): err = %q, want substring %q", tc.input, err.Error(), tc.wantError)
			}
		})
	}
}
