import type { CourseModule, ModuleStatus } from '@/types';

/**
 * Compute module status on the frontend based on lesson statuses and sort_order.
 *
 * Rules:
 * - completed: all lessons in module are completed
 * - available: first uncompleted module, or module whose previous module is fully completed
 * - locked: module after a not-yet-completed module
 */
export function computeModuleStatuses(
  modules: CourseModule[]
): Array<{ module: CourseModule; status: ModuleStatus }> {
  const sorted = [...modules].sort((a, b) => a.sort_order - b.sort_order);

  let previousCompleted = true;

  return sorted.map((mod) => {
    const allCompleted =
      mod.lessons.length > 0 && mod.lessons.every((l) => l.status === 'completed');
    const hasAnyProgress = mod.lessons.some(
      (l) => l.status === 'in_progress' || l.status === 'completed'
    );

    let status: ModuleStatus;

    if (allCompleted) {
      status = 'completed';
    } else if (previousCompleted || hasAnyProgress) {
      status = 'available';
    } else {
      status = 'locked';
    }

    // Only the previous module being fully completed unlocks next
    if (!allCompleted) {
      previousCompleted = false;
    }

    return { module: mod, status };
  });
}

/**
 * Calculate overall course progress as a percentage
 */
export function computeCourseProgress(modules: CourseModule[]): number {
  const allLessons = modules.flatMap((m) => m.lessons);
  if (allLessons.length === 0) return 0;
  const completed = allLessons.filter((l) => l.status === 'completed').length;
  return Math.round((completed / allLessons.length) * 100);
}
