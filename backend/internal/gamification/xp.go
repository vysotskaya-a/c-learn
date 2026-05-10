package gamification

// CalculateXP computes XP based on difficulty and attempt number.
// Penalty: -10% per failed attempt, minimum 30% of base.
func CalculateXP(difficulty string, attemptNumber int) int {
	base := map[string]int{
		"easy":   10,
		"medium": 25,
		"hard":   50,
	}

	xp, ok := base[difficulty]
	if !ok {
		xp = 10
	}

	penalty := float64(attemptNumber-1) * 0.1
	if penalty > 0.7 {
		penalty = 0.7
	}

	return int(float64(xp) * (1.0 - penalty))
}

// CalculateLevel returns level from total XP.
// Level thresholds: 0, 50, 150, 300, 500, 750, 1100, 1500, 2000, 2600, ...
func CalculateLevel(totalXP int) int {
	thresholds := []int{0, 50, 150, 300, 500, 750, 1100, 1500, 2000, 2600, 3300}
	level := 1
	for i, t := range thresholds {
		if totalXP >= t {
			level = i + 1
		} else {
			break
		}
	}
	return level
}
