import { useLeaderboard } from '@/hooks/useProfile';
import { useAuthStore } from '@/store/authStore';
import { CardSkeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import { Trophy, Zap, CheckCircle2 } from 'lucide-react';

export default function LeaderboardPage() {
  const { data, isLoading, error } = useLeaderboard();
  const { user } = useAuthStore();

  if (isLoading) {
    return (
      <div className="max-w-3xl mx-auto p-6 space-y-4">
        <CardSkeleton />
        <CardSkeleton />
        <CardSkeleton />
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-3xl mx-auto p-6">
        <EmptyState
          icon={<Trophy size={40} />}
          title="Рейтинг недоступен"
          description="Не удалось загрузить таблицу лидеров. Эта функция может быть временно недоступна."
        />
      </div>
    );
  }

  const entries = data?.entries || [];

  return (
    <div className="max-w-3xl mx-auto p-6 space-y-6 animate-fade-in">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-surface-900 dark:text-surface-100">Рейтинг</h1>
        {data?.current_user_rank && (
          <span className="text-sm text-surface-500">
            Ваша позиция: <strong className="text-brand-600 dark:text-brand-400">#{data.current_user_rank}</strong>
          </span>
        )}
      </div>

      {entries.length === 0 ? (
        <EmptyState
          icon={<Trophy size={40} />}
          title="Пока пусто"
          description="Рейтинг появится, когда кто-то решит первую задачу."
        />
      ) : (
        <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-700 overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="border-b border-surface-100 dark:border-surface-800">
                <th className="text-left px-5 py-3 text-xs font-semibold text-surface-500 uppercase tracking-wider w-16">
                  #
                </th>
                <th className="text-left px-5 py-3 text-xs font-semibold text-surface-500 uppercase tracking-wider">
                  Пользователь
                </th>
                <th className="text-right px-5 py-3 text-xs font-semibold text-surface-500 uppercase tracking-wider">
                  XP
                </th>
                <th className="text-right px-5 py-3 text-xs font-semibold text-surface-500 uppercase tracking-wider">
                  Решено
                </th>
              </tr>
            </thead>
            <tbody>
              {entries.map((entry) => {
                const isCurrentUser = entry.user_id === user?.id;
                return (
                  <tr
                    key={entry.user_id}
                    className={`
                      border-b border-surface-50 dark:border-surface-800/50 last:border-0
                      ${isCurrentUser ? 'bg-brand-50/50 dark:bg-brand-900/10' : ''}
                    `}
                  >
                    <td className="px-5 py-3.5">
                      {entry.rank <= 3 ? (
                        <span
                          className={`
                            w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold
                            ${entry.rank === 1 ? 'bg-amber-100 dark:bg-amber-900/30 text-amber-700' : ''}
                            ${entry.rank === 2 ? 'bg-gray-100 dark:bg-gray-800 text-gray-600' : ''}
                            ${entry.rank === 3 ? 'bg-orange-100 dark:bg-orange-900/30 text-orange-700' : ''}
                          `}
                        >
                          {entry.rank}
                        </span>
                      ) : (
                        <span className="text-sm text-surface-500 pl-1.5">{entry.rank}</span>
                      )}
                    </td>
                    <td className="px-5 py-3.5">
                      <span
                        className={`text-sm font-medium ${
                          isCurrentUser ? 'text-brand-700 dark:text-brand-400' : 'text-surface-800 dark:text-surface-200'
                        }`}
                      >
                        {isCurrentUser ? user?.username || 'Вы' : entry.username || `User #${entry.user_id.slice(0, 8)}`}
                      </span>
                    </td>
                    <td className="px-5 py-3.5 text-right">
                      <span className="inline-flex items-center gap-1 text-sm text-amber-600 dark:text-amber-400 font-medium">
                        <Zap size={14} />
                        {entry.total_xp}
                      </span>
                    </td>
                    <td className="px-5 py-3.5 text-right">
                      <span className="inline-flex items-center gap-1 text-sm text-surface-600 dark:text-surface-400">
                        <CheckCircle2 size={14} />
                        {entry.solved_count}
                      </span>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
