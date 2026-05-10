import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/Button';
import { FileQuestion } from 'lucide-react';

export default function NotFoundPage() {
  const navigate = useNavigate();

  return (
    <div className="flex items-center justify-center min-h-screen bg-surface-50 dark:bg-surface-950 p-8">
      <div className="text-center max-w-md">
        <div className="mx-auto w-20 h-20 rounded-2xl bg-surface-100 dark:bg-surface-800 flex items-center justify-center mb-6">
          <FileQuestion size={40} className="text-surface-400" />
        </div>
        <h1 className="text-6xl font-bold text-surface-900 dark:text-surface-100 mb-2">404</h1>
        <p className="text-lg text-surface-500 dark:text-surface-400 mb-8">
          Страница не найдена
        </p>
        <Button onClick={() => navigate('/dashboard')} size="lg">
          На главную
        </Button>
      </div>
    </div>
  );
}
