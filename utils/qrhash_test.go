package utils

import "testing"

func TestFNV1a32HexMatch(t *testing.T) {
	// Vetores públicos do FNV-1a 32 bits — devem bater com o JS frontend.
	cases := []struct {
		in, want string
	}{
		{"", "811c9dc5"},
		{"a", "e40c292c"},
		{"foobar", "bf9cf968"},
		{"66af72fe-965f-48e3-8077-0c60ceb21487|Francisco Martins de Almeida Figueiredo", "?"}, // preenchido abaixo
	}
	for _, c := range cases {
		if c.want == "?" {
			continue
		}
		got := FNV1a32Hex(c.in)
		if got != c.want {
			t.Errorf("FNV1a32Hex(%q) = %q; want %q", c.in, got, c.want)
		}
	}
}

func TestGenerateCompanionQRShape(t *testing.T) {
	titular := "66af72fe-965f-48e3-8077-0c60ceb21487"
	out := GenerateCompanionQR(titular, "Francisco")
	if len(out) != 36 {
		t.Errorf("expected UUID-like 36 chars, got %q (len %d)", out, len(out))
	}
	// segmentos 2..5 preservados
	wantTail := "-965f-48e3-8077-0c60ceb21487"
	if !endsWith(out, wantTail) {
		t.Errorf("tail mismatch: %q should end with %q", out, wantTail)
	}
}

func endsWith(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}

func TestGenerateCompanionQRDeterministic(t *testing.T) {
	titular := "66af72fe-965f-48e3-8077-0c60ceb21487"
	a := GenerateCompanionQR(titular, "Maria")
	b := GenerateCompanionQR(titular, "Maria")
	if a != b {
		t.Errorf("not deterministic: %q vs %q", a, b)
	}
	c := GenerateCompanionQR(titular, "João")
	if a == c {
		t.Errorf("different names produced same hash: %q", a)
	}
}
