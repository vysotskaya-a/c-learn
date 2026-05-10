import { ReactNode } from 'react';

interface AuthLayoutProps {
  children: ReactNode;
}

export function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="min-h-screen flex bg-surface-50 dark:bg-surface-950">
      {/* Left side - branding */}
      <div className="hidden lg:flex lg:w-1/2 bg-brand-600 relative overflow-hidden items-center justify-center">
        <div className="absolute inset-0 opacity-10">
          <div className="absolute top-10 left-10 text-[200px] font-mono font-bold text-white/20 select-none">
            {'{'}
          </div>
          <div className="absolute bottom-10 right-10 text-[200px] font-mono font-bold text-white/20 select-none">
            {'}'}
          </div>
          <div className="absolute top-1/3 left-1/4 text-6xl font-mono text-white/10 select-none rotate-12">
            #include &lt;stdio.h&gt;
          </div>
          <div className="absolute bottom-1/3 right-1/4 text-5xl font-mono text-white/10 select-none -rotate-6">
            int main()
          </div>
        </div>
        <div className="relative z-10 text-center px-12">
          <div className="w-20 h-20 rounded-2xl bg-white/20 backdrop-blur flex items-center justify-center mx-auto mb-8">
            <span className="text-white font-mono font-bold text-4xl">C</span>
          </div>
          <h1 className="text-4xl font-display font-bold text-white mb-4">C-Learn</h1>
          <p className="text-lg text-white/80 max-w-md">
            Изучай язык C через практику. Решай задачи, получай XP, открывай достижения.
          </p>
        </div>
      </div>

      {/* Right side - form */}
      <div className="flex-1 flex items-center justify-center p-6">
        <div className="w-full max-w-md">{children}</div>
      </div>
    </div>
  );
}
