import { apiClient } from './client';
import type {
  RegisterRequest,
  RegisterResponse,
  LoginRequest,
  TokenResponse,
  UserInfo,
} from '@/types';

export const authApi = {
  register: (data: RegisterRequest) =>
    apiClient.post<RegisterResponse>('/api/v1/auth/register', data),

  login: (data: LoginRequest) =>
    apiClient.post<TokenResponse>('/api/v1/auth/login', data),

  refresh: (refreshToken: string) =>
    apiClient.post<TokenResponse>('/api/v1/auth/refresh', { refresh_token: refreshToken }),

  me: () =>
    apiClient.get<UserInfo>('/api/v1/auth/me'),
};
