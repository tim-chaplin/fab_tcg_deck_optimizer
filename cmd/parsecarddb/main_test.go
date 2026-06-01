package main

import "testing"

func set(words ...string) map[string]bool {
	m := map[string]bool{}
	for _, w := range words {
		m[w] = true
	}
	return m
}

// TestPoolFilterMatches covers the class/talent subset semantics that express a hero's pool.
func TestPoolFilterMatches(t *testing.T) {
	card := func(types string) Card { return Card{Name: "X", Types: types, SilverAgeLegal: "Yes"} }

	viserai := poolFilter{classes: set("Runeblade"), silverAge: true} // no talent
	aurora := poolFilter{classes: set("Runeblade"), talents: set("Lightning"), silverAge: true}

	cases := []struct {
		name   string
		filter poolFilter
		card   Card
		want   bool
	}{
		{"viserai: runeblade", viserai, card("Runeblade, Action, Attack"), true},
		{"viserai: generic always legal", viserai, card("Generic, Action"), true},
		{"viserai: classless no-talent", viserai, card("Action, Attack"), true},
		{"viserai: rejects lightning", viserai, card("Lightning, Runeblade, Action"), false},
		{"viserai: rejects off-class", viserai, card("Warrior, Action, Attack"), false},
		{"aurora: lightning runeblade", aurora, card("Lightning, Runeblade, Action"), true},
		{"aurora: generic lightning", aurora, card("Lightning, Generic, Instant"), true},
		{"aurora: plain runeblade", aurora, card("Runeblade, Action"), true},
		{"aurora: rejects shadow", aurora, card("Shadow, Runeblade, Aura"), false},
		// require-talents narrows a pool to just the talented cards.
		{"require lightning: has it", poolFilter{classes: set("Runeblade"), talents: set("Lightning"), require: set("Lightning")}, card("Lightning, Runeblade, Action"), true},
		{"require lightning: no-talent excluded", poolFilter{classes: set("Runeblade"), talents: set("Lightning"), require: set("Lightning")}, card("Runeblade, Action"), false},
		// any-talent ignores the subset check.
		{"any-talent admits shadow", poolFilter{anyTalent: true}, card("Shadow, Runeblade, Aura"), true},
		// modeled keeps weapons (the optimizer models them) but drops unmodeled types.
		{"modeled keeps weapon", poolFilter{modeled: true}, card("Runeblade, Weapon, Sword"), true},
		{"modeled drops equipment", poolFilter{modeled: true}, card("Generic, Equipment, Arms"), false},
		// silver-age gate.
		{"silver-age drops illegal", poolFilter{silverAge: true}, Card{Name: "X", Types: "Generic, Action", SilverAgeLegal: "No"}, false},
		{"empty classes = any class", poolFilter{}, card("Warrior, Action"), true},
	}
	for _, c := range cases {
		if got := c.filter.matches(c.card); got != c.want {
			t.Errorf("%s: matches = %v, want %v", c.name, got, c.want)
		}
	}
}

// TestPoolFilterExcludeBanned checks the cross-reference against our internal/format banlist.
func TestPoolFilterExcludeBanned(t *testing.T) {
	f := poolFilter{excludeBanned: true}
	if f.matches(Card{Name: "Ball Lightning", Types: "Lightning, Action, Attack", SilverAgeLegal: "Yes"}) {
		t.Error("Ball Lightning (on our banlist) should be excluded by -exclude-banned")
	}
	if !f.matches(Card{Name: "Aether Slash", Types: "Runeblade, Action, Attack", SilverAgeLegal: "Yes"}) {
		t.Error("Aether Slash (not banned) should pass -exclude-banned")
	}
}
