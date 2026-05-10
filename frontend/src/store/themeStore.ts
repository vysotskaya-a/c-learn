import { create } from 'zustand';

interface ThemeState {
  dark: boolean;
  toggle: () => void;
}

const getInitialTheme = (): boolean => {
  try {
    const stored = localStorage.getItem('c-learn-theme');
    if (stored !== null) return stored === 'dark';
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  } catch {
    return false;
  }
};

export const useThemeStore = create<ThemeState>((set, get) => ({
  dark: getInitialTheme(),
  toggle: () => {
    const next = !get().dark;
    set({ dark: next });
    localStorage.setItem('c-learn-theme', next ? 'dark' : 'light');
    if (next) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  },
}));

// Apply on load
if (getInitialTheme()) {
  document.documentElement.classList.add('dark');
}
