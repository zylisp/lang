package sexpr

import "testing"

func TestNumberString(t *testing.T) {
	tests := []struct {
		value    int64
		expected string
	}{
		{42, "42"},
		{-17, "-17"},
		{0, "0"},
	}

	for _, tt := range tests {
		n := Number{Value: tt.value}
		if got := n.String(); got != tt.expected {
			t.Errorf("Number(%d).String() = %q, want %q",
				tt.value, got, tt.expected)
		}
	}
}

func TestSymbolString(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"x", "x"},
		{"+", "+"},
		{"lambda", "lambda"},
	}

	for _, tt := range tests {
		s := Symbol{Name: tt.name}
		if got := s.String(); got != tt.expected {
			t.Errorf("Symbol(%q).String() = %q, want %q",
				tt.name, got, tt.expected)
		}
	}
}

func TestBoolString(t *testing.T) {
	tests := []struct {
		value    bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	}

	for _, tt := range tests {
		b := Bool{Value: tt.value}
		if got := b.String(); got != tt.expected {
			t.Errorf("Bool(%v).String() = %q, want %q",
				tt.value, got, tt.expected)
		}
	}
}

func TestListString(t *testing.T) {
	tests := []struct {
		name     string
		list     List
		expected string
	}{
		{
			"empty list",
			List{Elements: []SExpr{}},
			"()",
		},
		{
			"single element",
			List{Elements: []SExpr{Number{Value: 42}}},
			"(42)",
		},
		{
			"multiple elements",
			List{Elements: []SExpr{
				Symbol{Name: "+"},
				Number{Value: 1},
				Number{Value: 2},
			}},
			"(+ 1 2)",
		},
		{
			"nested list",
			List{Elements: []SExpr{
				Symbol{Name: "+"},
				List{Elements: []SExpr{
					Symbol{Name: "*"},
					Number{Value: 2},
					Number{Value: 3},
				}},
				Number{Value: 4},
			}},
			"(+ (* 2 3) 4)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.String(); got != tt.expected {
				t.Errorf("List.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}
