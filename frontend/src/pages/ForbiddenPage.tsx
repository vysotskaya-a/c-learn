import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/Button';
import { ShieldX } from 'lucide-react';

export default function ForbiddenPage() {
  const navigate = useNavigate();

  return (
    <div className="flex items-center justify-center min-h-screen bg-surface-50 dark:bg-surface-950 p-8">
      <div className="text-center max-w-md">
        <div className="mx-auto w-20 h-20 rounded-2xl bg-red-100 dark:bg-red-900/20 flex items-center justify-center mb-6">
          <ShieldX size={40} className="text-red-500" />
        </div>
        <h1 className="text-6xl font-bold text-surface-900 dark:text-surface-100 mb-2">403</h1>
        <p className="text-lg text-surface-500 dark:text-surface-400 mb-8">
          Доступ запрещён
        </p>
        <Button onClick={() => navigate('/dashboard')} size="lg">
          На главную
        </Button>
      </div>
    </div>
  );
}
