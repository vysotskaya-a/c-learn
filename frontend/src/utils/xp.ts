/**
 * Calculate XP needed for a given level.
 * Simple formula: each level requires level * 100 XP.
 */
export function xpForLevel(level: number): number {
  return level * 100;
}

/**
 * Calculate total XP needed from 0 to reach a given level.
 */
export function totalXpForLevel(level: number): number {
  return (level * (level + 1)) / 2 * 100;
}

/**
 * Get progress within current level as a percentage (0-100).
 */
export function getLevelProgress(totalXp: number, level: number): number {
  if (level <= 0) return 0;
  const xpForCurrent = totalXpForLevel(level - 1);
  const xpForNext = totalXpForLevel(level);
  const range = xpForNext - xpForCurrent;
  if (range <= 0) return 100;
  const progress = ((totalXp - xpForCurrent) / range) * 100;
  return Math.min(100, Math.max(0, Math.round(progress)));
}
