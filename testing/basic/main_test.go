package main

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"pos", 1, 2, 3},
		{"zero", 0, 5, 5},
		{"neg", -2, -3, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.a, tt.b); got != tt.expected {
				t.Fatalf("Add(%d,%d)=%d, want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		got, err := Div(10, 2)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if got != 5 {
			t.Fatalf("Div(10,2)=%d, want 5", got)
		}
	})

	t.Run("divide_by_zero", func(t *testing.T) {
		_, err := Div(10, 0)
		if err == nil {
			t.Fatalf("expected err")
		}
	})
}

func TestNormalizeName(t *testing.T) {
	got := NormalizeName("  Reo  ")
	if got != "reo" {
		t.Fatalf("NormalizeName()=%q, want %q", got, "reo")
	}
}

func TestIsEven(t *testing.T) {
	if !IsEven(2) {
		t.Fatalf("expected true")
	}
	if IsEven(3) {
		t.Fatalf("expected false")
	}
}
