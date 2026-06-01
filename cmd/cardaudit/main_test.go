package main

import "testing"

func TestAsStr(t *testing.T) {
	cases := []struct {
		in   any
		want string
	}{
		{"  0 ", "0"},     // string, trimmed
		{float64(4), "4"}, // JSON number → integer form
		{nil, ""},         // null
		{float64(0), "0"},
	}
	for _, c := range cases {
		if got := asStr(c.in); got != c.want {
			t.Errorf("asStr(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestParseInt(t *testing.T) {
	if n, ok := parseInt(" 3 "); !ok || n != 3 {
		t.Errorf(`parseInt(" 3 ") = %d,%v, want 3,true`, n, ok)
	}
	if _, ok := parseInt(""); ok {
		t.Error("parseInt(\"\") should be ok=false (no value to compare)")
	}
	if _, ok := parseInt("X"); ok {
		t.Error("parseInt(\"X\") should be ok=false (non-numeric)")
	}
}

func TestNormName(t *testing.T) {
	if got := normName("Burn Up||Shock"); got != "Burn Up // Shock" {
		t.Errorf("normName(||) = %q, want %q", got, "Burn Up // Shock")
	}
	if got := normName("Aether Slash"); got != "Aether Slash" {
		t.Errorf("normName(plain) = %q, want unchanged", got)
	}
}
