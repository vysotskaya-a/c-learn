import { create } from 'zustand';
import type { UserInfo } from '@/types';

interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  user: UserInfo | null;
  expiresAt: number | null;

  setTokens: (access: string, refresh: string, expiresIn: number) => void;
  setUser: (user: UserInfo) => void;
  logout: () => void;
  isAuthenticated: () => boolean;
  isAdmin: () => boolean;
}

const STORAGE_KEY = 'c-learn-auth';

function loadPersistedState() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return {};
    const data = JSON.parse(raw);
    return {
      accessToken: data.accessToken ?? null,
      refreshToken: data.refreshToken ?? null,
      user: data.user ?? null,
      expiresAt: data.expiresAt ?? null,
    };
  } catch {
    return {};
  }
}

function persistState(state: Partial<AuthState>) {
  try {
    localStorage.setItem(
      STORAGE_KEY,
      JSON.stringify({
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        user: state.user,
        expiresAt: state.expiresAt,
      })
    );
  } catch {
    // ignore
  }
}

export const useAuthStore = create<AuthState>((set, get) => ({
  accessToken: null,
  refreshToken: null,
  user: null,
  expiresAt: null,
  ...loadPersistedState(),

  setTokens: (access, refresh, expiresIn) => {
    const expiresAt = Date.now() + expiresIn * 1000;
    set({ accessToken: access, refreshToken: refresh, expiresAt });
    persistState({ ...get(), accessToken: access, refreshToken: refresh, expiresAt });
  },

  setUser: (user) => {
    set({ user });
    persistState({ ...get(), user });
  },

  logout: () => {
    set({ accessToken: null, refreshToken: null, user: null, expiresAt: null });
    localStorage.removeItem(STORAGE_KEY);
  },

  isAuthenticated: () => {
    const { accessToken } = get();
    return !!accessToken;
  },

  isAdmin: () => {
    const { user } = get();
    return user?.role === 'admin';
  },
}));
