import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCourseTree } from '@/hooks/useCourseTree';
import { useProfile } from '@/hooks/useProfile';
import { computeModuleStatuses, computeCourseProgress } from '@/utils/moduleStatus';
import { getLevelProgress } from '@/utils/xp';
import { CardSkeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import { Button } from '@/components/ui/Button';
import {
  ChevronDown,
  ChevronRight,
  CheckCircle2,
  Circle,
  Lock,
  BookOpen,
  Zap,
  Star,
  Trophy,
} from 'lucide-react';
import type { CourseModule, ModuleStatus, LessonStatus } from '@/types';

function StatusBadge({ status }: { status: ModuleStatus }) {
  if (status === 'completed')
    return (
      <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400">
        <CheckCircle2 size={12} /> Завершён
      </span>
    );
  if (status === 'available')
    return (
      <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-brand-100 dark:bg-brand-900/30 text-brand-700 dark:text-brand-400">
        <Circle size={12} /> Доступен
      </span>
    );
  return (
    <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-surface-100 dark:bg-surface-800 text-surface-500 dark:text-surface-400">
      <Lock size={12} /> Заблокирован
    </span>
  );
}

function LessonStatusIcon({ status }: { status: LessonStatus }) {
  if (status === 'completed') return <CheckCircle2 size={16} className="text-emerald-500" />;
  if (status === 'in_progress') return <Circle size={16} className="text-brand-500" />;
  return <Lock size={16} className="text-surface-300 dark:text-surface-600" />;
}

function ModuleCard({
  module,
  status,
  defaultOpen,
}: {
  module: CourseModule;
  status: ModuleStatus;
  defaultOpen: boolean;
}) {
  const [open, setOpen] = useState(defaultOpen);
  const navigate = useNavigate();
  const sortedLessons = [...module.lessons].sort((a, b) => a.sort_order - b.sort_order);
  const completedCount = module.lessons.filter((l) => l.status === 'completed').length;
  const total = module.lessons.length;
  const progressPct = total > 0 ? Math.round((completedCount / total) * 100) : 0;

  return (
    <div
      className={`
        rounded-xl border transition-all
        ${status === 'locked'
          ? 'border-surface-100 dark:border-surface-800 opacity-60'
          : 'border-surface-200 dark:border-surface-700 bg-white dark:bg-surface-900'
        }
      `}
    >
      <button
        onClick={() => status !== 'locked' && setOpen(!open)}
        disabled={status === 'locked'}
        className="w-full flex items-center gap-4 p-5 text-left"
      >
        <div
          className={`
            w-10 h-10 rounded-lg flex items-center justify-center shrink-0
            ${status === 'completed' ? 'bg-emerald-100 dark:bg-emerald-900/30' : ''}
            ${status === 'available' ? 'bg-brand-100 dark:bg-brand-900/30' : ''}
            ${status === 'locked' ? 'bg-surface-100 dark:bg-surface-800' : ''}
          `}
        >
          <BookOpen
            size={20}
            className={`
              ${status === 'completed' ? 'text-emerald-600 dark:text-emerald-400' : ''}
              ${status === 'available' ? 'text-brand-600 dark:text-brand-400' : ''}
              ${status === 'locked' ? 'text-surface-400' : ''}
            `}
          />
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <h3 className="font-semibold text-surface-900 dark:text-surface-100 truncate">{module.title}</h3>
            <StatusBadge status={status} />
          </div>
          <div className="flex items-center gap-3">
            <span className="text-xs text-surface-500">
              {completedCount}/{total} уроков
            </span>
            <div className="flex-1 max-w-[200px] h-1.5 bg-surface-100 dark:bg-surface-800 rounded-full overflow-hidden">
              <div
                className="h-full bg-emerald-500 rounded-full transition-all duration-500"
                style={{ width: `${progressPct}%` }}
              />
            </div>
          </div>
        </div>

        {status !== 'locked' && (
          <div className="shrink-0 text-surface-400">
            {open ? <ChevronDown size={20} /> : <ChevronRight size={20} />}
          </div>
        )}
      </button>

      {open && status !== 'locked' && (
        <div className="border-t border-surface-100 dark:border-surface-800 px-5 pb-4">
          {sortedLessons.map((lesson) => (
            <div
              key={lesson.id}
              className="flex items-center gap-3 py-3 border-b border-surface-50 dark:border-surface-800/50 last:border-0"
            >
              <LessonStatusIcon status={lesson.status} />
              <span
                className={`flex-1 text-sm ${
                  lesson.status === 'completed'
                    ? 'text-surface-500 dark:text-surface-400'
                    : 'text-surface-800 dark:text-surface-200'
                }`}
              >
                {lesson.title}
              </span>
              {lesson.status !== 'not_started' || status === 'available' ? (
                <Button
                  size="sm"
                  variant={lesson.status === 'completed' ? 'ghost' : 'primary'}
                  onClick={() => navigate(`/lessons/${lesson.id}`)}
                >
                  {lesson.status === 'completed'
                    ? 'Повторить'
                    : lesson.status === 'in_progress'
                      ? 'Продолжить'
                      : 'Начать'}
                </Button>
              ) : null}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default function DashboardPage() {
  const { data: tree, isLoading: treeLoading, error: treeError } = useCourseTree();
  const { data: profile } = useProfile();

  if (treeLoading) {
    return (
      <div className="max-w-4xl mx-auto p-6 space-y-4">
        <CardSkeleton />
        <CardSkeleton />
        <CardSkeleton />
      </div>
    );
  }

  if (treeError) {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <EmptyState title="Ошибка загрузки" description="Не удалось загрузить дерево курса. Попробуйте позже." />
      </div>
    );
  }

  const modules = tree?.modules || [];
  const withStatuses = computeModuleStatuses(modules);
  const progress = computeCourseProgress(modules);
  const firstAvailableIdx = withStatuses.findIndex((m) => m.status === 'available');

  return (
    <div className="max-w-4xl mx-auto p-6 space-y-6 animate-fade-in">
      {/* Stats bar */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-700 p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-9 h-9 rounded-lg bg-brand-100 dark:bg-brand-900/30 flex items-center justify-center">
              <BookOpen size={18} className="text-brand-600 dark:text-brand-400" />
            </div>
            <span className="text-sm font-medium text-surface-500">Прогресс курса</span>
          </div>
          <div className="flex items-end gap-2">
            <span className="text-3xl font-bold text-surface-900 dark:text-surface-100">{progress}%</span>
          </div>
          <div className="mt-2 h-2 bg-surface-100 dark:bg-surface-800 rounded-full overflow-hidden">
            <div
              className="h-full bg-brand-500 rounded-full transition-all duration-700"
              style={{ width: `${progress}%` }}
            />
          </div>
        </div>

        {profile && (
          <>
            <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-700 p-5">
              <div className="flex items-center gap-3 mb-3">
                <div className="w-9 h-9 rounded-lg bg-amber-100 dark:bg-amber-900/30 flex items-center justify-center">
                  <Zap size={18} className="text-amber-600 dark:text-amber-400" />
                </div>
                <span className="text-sm font-medium text-surface-500">XP / Уровень</span>
              </div>
              <div className="flex items-end gap-2">
                <span className="text-3xl font-bold text-surface-900 dark:text-surface-100">
                  {profile.total_xp}
                </span>
                <span className="text-sm text-surface-400 mb-1">XP</span>
              </div>
              <div className="mt-2 flex items-center gap-2">
                <Star size={14} className="text-amber-500" />
                <span className="text-sm text-surface-600 dark:text-surface-400">
                  Уровень {profile.level}
                </span>
                <div className="flex-1 h-1.5 bg-surface-100 dark:bg-surface-800 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-amber-500 rounded-full transition-all"
                    style={{ width: `${getLevelProgress(profile.total_xp, profile.level)}%` }}
                  />
                </div>
              </div>
            </div>

            <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-700 p-5">
              <div className="flex items-center gap-3 mb-3">
                <div className="w-9 h-9 rounded-lg bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center">
                  <Trophy size={18} className="text-emerald-600 dark:text-emerald-400" />
                </div>
                <span className="text-sm font-medium text-surface-500">Решено задач</span>
              </div>
              <span className="text-3xl font-bold text-surface-900 dark:text-surface-100">
                {profile.solved_count}
              </span>
            </div>
          </>
        )}
      </div>

      {/* Modules */}
      <div>
        <h2 className="text-xl font-bold text-surface-900 dark:text-surface-100 mb-4">Модули курса</h2>
        {modules.length === 0 ? (
          <EmptyState title="Курс пуст" description="Модули ещё не добавлены. Обратитесь к администратору." />
        ) : (
          <div className="space-y-3">
            {withStatuses.map(({ module, status }, idx) => (
              <ModuleCard
                key={module.id}
                module={module}
                status={status}
                defaultOpen={idx === firstAvailableIdx}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
