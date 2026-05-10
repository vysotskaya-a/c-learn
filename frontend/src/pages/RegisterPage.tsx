import { useState, FormEvent } from 'react';
import { Link } from 'react-router-dom';
import { AuthLayout } from '@/components/layout/AuthLayout';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { useRegister } from '@/hooks/useAuth';

export default function RegisterPage() {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const register = useRegister();

  const validate = (): boolean => {
    const e: Record<string, string> = {};
    if (!username.trim() || username.trim().length < 3) e.username = 'Nickname минимум 3 символа';
    if (!email.trim()) e.email = 'Email обязателен';
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) e.email = 'Некорректный формат email';
    if (!password || password.length < 8) e.password = 'Пароль минимум 8 символов';
    if (password !== confirmPassword) e.confirmPassword = 'Пароли не совпадают';
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleSubmit = (ev: FormEvent) => {
    ev.preventDefault();
    if (!validate()) return;
    register.mutate({ username: username.trim(), email: email.trim(), password });
  };

  return (
    <AuthLayout>
      <div className="lg:hidden flex items-center gap-3 mb-8">
        <div className="w-10 h-10 rounded-xl bg-brand-600 flex items-center justify-center">
          <span className="text-white font-mono font-bold text-xl">C</span>
        </div>
        <span className="font-display font-bold text-xl text-surface-900 dark:text-surface-100">C-Learn</span>
      </div>

      <h2 className="text-2xl font-bold text-surface-900 dark:text-surface-100 mb-1">Регистрация</h2>
      <p className="text-surface-500 dark:text-surface-400 mb-8">Создайте аккаунт для начала обучения</p>

      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Nickname"
          placeholder="student1"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          error={errors.username}
          autoComplete="username"
        />
        <Input
          label="Email"
          type="email"
          placeholder="student@example.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          error={errors.email}
          autoComplete="email"
        />
        <Input
          label="Пароль"
          type="password"
          placeholder="Минимум 8 символов"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          error={errors.password}
          autoComplete="new-password"
        />
        <Input
          label="Подтверждение пароля"
          type="password"
          placeholder="Повторите пароль"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          error={errors.confirmPassword}
          autoComplete="new-password"
        />
        <Button type="submit" className="w-full" size="lg" loading={register.isPending}>
          Зарегистрироваться
        </Button>
      </form>

      <p className="text-center text-sm text-surface-500 dark:text-surface-400 mt-6">
        Уже есть аккаунт?{' '}
        <Link to="/login" className="text-brand-600 dark:text-brand-400 font-medium hover:underline">
          Войти
        </Link>
      </p>
    </AuthLayout>
  );
}
