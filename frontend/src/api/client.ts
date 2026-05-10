import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/store/authStore';
import type { TokenResponse, ApiError } from '@/types';

/**
 * BASE_URL resolution:
 *
 * Option A (dev with Vite proxy — recommended):
 *   VITE_API_BASE_URL is empty or not set → requests go to same origin,
 *   Vite proxy rewrites /api/* → http://localhost:8080/api/*
 *
 * Option B (direct to gateway — production or no proxy):
 *   VITE_API_BASE_URL=http://localhost:8080
 *   Gateway MUST have CORS configured for the frontend origin.
 *
 * In BOTH cases the endpoint paths in code are absolute: "/api/v1/auth/login"
 * so the BASE_URL must NOT contain "/api/v1" — that would double the prefix.
 */
const RAW_BASE = import.meta.env.VITE_API_BASE_URL || '';
// Strip trailing slashes and strip /api/v1 suffix if someone accidentally added it
const BASE_URL = RAW_BASE
  .replace(/\/+$/, '')
  .replace(/\/api\/v1\/?$/, '');

export const apiClient = axios.create({
  baseURL: BASE_URL,
  headers: { 'Content-Type': 'application/json' },
  timeout: 30000,
});

// --- Request interceptor: attach Bearer token ---
apiClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().accessToken;
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// --- Refresh queue to prevent refresh storm ---
let isRefreshing = false;
let failedQueue: {
  resolve: (token: string) => void;
  reject: (err: unknown) => void;
}[] = [];

function processQueue(error: unknown, token: string | null) {
  failedQueue.forEach((p) => {
    if (error) {
      p.reject(error);
    } else {
      p.resolve(token!);
    }
  });
  failedQueue = [];
}

// --- Response interceptor: handle 401 + refresh ---
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiError>) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    // Don't retry refresh/login/register endpoints
    const url = originalRequest?.url || '';
    const isAuthEndpoint =
      url.includes('/auth/login') ||
      url.includes('/auth/register') ||
      url.includes('/auth/refresh');

    if (error.response?.status === 401 && !originalRequest._retry && !isAuthEndpoint) {
      if (isRefreshing) {
        // Queue this request until refresh completes
        return new Promise((resolve, reject) => {
          failedQueue.push({
            resolve: (token: string) => {
              if (originalRequest.headers) {
                originalRequest.headers.Authorization = `Bearer ${token}`;
              }
              resolve(apiClient(originalRequest));
            },
            reject,
          });
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const refreshToken = useAuthStore.getState().refreshToken;
      if (!refreshToken) {
        useAuthStore.getState().logout();
        window.location.href = '/login';
        return Promise.reject(error);
      }

      try {
        // Use raw axios (not apiClient) to avoid interceptor loop,
        // but with the same BASE_URL
        const { data } = await axios.post<TokenResponse>(
          `${BASE_URL}/api/v1/auth/refresh`,
          { refresh_token: refreshToken },
          { headers: { 'Content-Type': 'application/json' } }
        );

        useAuthStore.getState().setTokens(data.access_token, data.refresh_token, data.expires_in);
        processQueue(null, data.access_token);

        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${data.access_token}`;
        }
        return apiClient(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        useAuthStore.getState().logout();
        window.location.href = '/login';
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);

/**
 * Extract user-friendly error message from API error
 */
export function getErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as ApiError | undefined;
    if (data?.message) return data.message;
    if (error.response?.status === 401) return 'Неверный email или пароль';
    if (error.response?.status === 403) return 'Доступ запрещён';
    if (error.response?.status === 404) return 'Ресурс не найден';
    if (error.response?.status === 409) return 'Email или username уже заняты';
    if (error.response?.status === 500) return 'Внутренняя ошибка сервера';
    if (error.message) return error.message;
  }
  return 'Произошла неизвестная ошибка';
}
