package runner

import "testing"

func TestNormalizeOutput(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "42", want: "42"},
		{in: "42\n", want: "42"},
		{in: "42\n\n", want: "42"},
		{in: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := NormalizeOutput(tt.in); got != tt.want {
				t.Fatalf("NormalizeOutput(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
