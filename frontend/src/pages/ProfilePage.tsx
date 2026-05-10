import { useFetchUser } from '@/hooks/useAuth';
import { useProfile } from '@/hooks/useProfile';
import { getInitials, getAvatarColor } from '@/utils/avatar';
import { getLevelProgress } from '@/utils/xp';
import { ProfileSkeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import { Zap, Star, CheckCircle2, Award, Calendar } from 'lucide-react';
import type { ProfileAchievement } from '@/types';

function AchievementCard({ achievement }: { achievement: ProfileAchievement }) {
  return (
    <div className="flex items-start gap-3 p-4 rounded-xl bg-white dark:bg-surface-900 border border-surface-200 dark:border-surface-700">
      <div className="w-10 h-10 rounded-lg bg-amber-100 dark:bg-amber-900/30 flex items-center justify-center shrink-0">
        <Award size={20} className="text-amber-600 dark:text-amber-400" />
      </div>
      <div className="min-w-0">
        <p className="text-sm font-semibold text-surface-900 dark:text-surface-100">{achievement.title}</p>
        <p className="text-xs text-surface-500 dark:text-surface-400 mt-0.5">{achievement.description}</p>
        <p className="text-xs text-surface-400 mt-1">
          {new Date(achievement.awarded_at).toLocaleDateString('ru-RU')}
        </p>
      </div>
    </div>
  );
}

export default function ProfilePage() {
  const { data: user, isLoading: userLoading } = useFetchUser();
  const { data: profile, isLoading: profileLoading, error: profileError } = useProfile();

  if (userLoading || profileLoading) {
    return (
      <div className="max-w-3xl mx-auto p-6">
        <ProfileSkeleton />
      </div>
    );
  }

  if (!user) {
    return (
      <div className="max-w-3xl mx-auto p-6">
        <EmptyState title="Ошибка" description="Не удалось загрузить данные пользователя." />
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto p-6 space-y-6 animate-fade-in">
      {/* User card */}
      <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-700 p-6">
        <div className="flex items-center gap-5">
          <div
            className={`w-20 h-20 rounded-full flex items-center justify-center text-white text-2xl font-bold ${getAvatarColor(user.username)}`}
          >
            {getInitials(user.username)}
          </div>
          <div>
            <h1 className="text-2xl font-bold text-surface-900 dark:text-surface-100">{user.username}</h1>
            <p className="text-sm text-surface-500">{user.email}</p>
            <div className="flex items-center gap-1.5 mt-1 text-xs text-surface-400">
              <Calendar size={12} />
              Зарегистрирован {new Date(user.created_at).toLocaleDateString('ru-RU')}
            </div>
          </div>
        </div>
      </div>

      {/* Stats */}
      {profile && (
        <>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-700 p-5 text-center">
              <Zap size={24} className="mx-auto text-amber-500 mb-2" />
              <p className="text-3xl font-bold text-surface-900 dark:text-surface-100">{profile.total_xp}</p>
              <p className="text-xs text-surface-500 mt-1">Всего XP</p>
            </div>
            <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-700 p-5 text-center">
              <Star size={24} className="mx-auto text-brand-500 mb-2" />
              <p className="text-3xl font-bold text-surface-900 dark:text-surface-100">{profile.level}</p>
              <p className="text-xs text-surface-500 mt-1">Уровень</p>
              <div className="mt-2 h-2 bg-surface-100 dark:bg-surface-800 rounded-full overflow-hidden">
                <div
                  className="h-full bg-brand-500 rounded-full transition-all"
                  style={{ width: `${getLevelProgress(profile.total_xp, profile.level)}%` }}
                />
              </div>
            </div>
            <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-700 p-5 text-center">
              <CheckCircle2 size={24} className="mx-auto text-emerald-500 mb-2" />
              <p className="text-3xl font-bold text-surface-900 dark:text-surface-100">{profile.solved_count}</p>
              <p className="text-xs text-surface-500 mt-1">Решено задач</p>
            </div>
          </div>

          {/* Achievements */}
          <div>
            <h2 className="text-lg font-bold text-surface-900 dark:text-surface-100 mb-4">Достижения</h2>
            {profile.achievements.length === 0 ? (
              <EmptyState
                icon={<Award size={40} />}
                title="Пока нет достижений"
                description="Решайте задачи, чтобы получить первые бейджи!"
              />
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                {profile.achievements.map((a) => (
                  <AchievementCard key={a.code} achievement={a} />
                ))}
              </div>
            )}
          </div>
        </>
      )}

      {profileError && (
        <EmptyState
          title="Данные геймификации недоступны"
          description="Не удалось загрузить XP и достижения. Попробуйте позже."
        />
      )}
    </div>
  );
}
