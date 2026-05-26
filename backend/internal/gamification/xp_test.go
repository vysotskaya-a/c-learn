package gamification

import "testing"

func TestCalculateXP(t *testing.T) {
	tests := []struct {
		name          string
		difficulty    string
		attemptNumber int
		want          int
	}{
		// PMI: base × max(0.3, 1 − 0.1 × failed_attempts)
		{name: "easy first try", difficulty: "easy", attemptNumber: 1, want: 10},
		{name: "medium first try", difficulty: "medium", attemptNumber: 1, want: 25},
		{name: "hard first try", difficulty: "hard", attemptNumber: 1, want: 50},
		{name: "easy second try", difficulty: "easy", attemptNumber: 2, want: 9},
		{name: "medium third try", difficulty: "medium", attemptNumber: 3, want: 20},
		{name: "hard max penalty", difficulty: "hard", attemptNumber: 10, want: 15},
		{name: "unknown difficulty defaults easy", difficulty: "unknown", attemptNumber: 1, want: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateXP(tt.difficulty, tt.attemptNumber)
			if got != tt.want {
				t.Fatalf("CalculateXP(%q, %d) = %d, want %d", tt.difficulty, tt.attemptNumber, got, tt.want)
			}
		})
	}
}

func TestCalculateLevel(t *testing.T) {
	tests := []struct {
		totalXP int
		want    int
	}{
		{totalXP: 0, want: 1},
		{totalXP: 49, want: 1},
		{totalXP: 50, want: 2},
		{totalXP: 149, want: 2},
		{totalXP: 150, want: 3},
		{totalXP: 3300, want: 11},
		{totalXP: 10000, want: 11},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := CalculateLevel(tt.totalXP); got != tt.want {
				t.Fatalf("CalculateLevel(%d) = %d, want %d", tt.totalXP, got, tt.want)
			}
		})
	}
}
