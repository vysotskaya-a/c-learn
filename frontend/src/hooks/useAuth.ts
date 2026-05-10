import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { authApi } from '@/api/auth';
import { useAuthStore } from '@/store/authStore';
import { useToastStore } from '@/store/toastStore';
import { getErrorMessage } from '@/api/client';
import { useNavigate } from 'react-router-dom';
import type { LoginRequest, RegisterRequest } from '@/types';

export function useLogin() {
  const { setTokens } = useAuthStore();
  const { addToast } = useToastStore();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: LoginRequest) => authApi.login(data),
    onSuccess: async (res) => {
      queryClient.clear();
      setTokens(res.data.access_token, res.data.refresh_token, res.data.expires_in);
      // Fetch user info after login
      try {
        const meRes = await authApi.me();
        useAuthStore.getState().setUser(meRes.data);
      } catch {
        // non-critical
      }
      addToast('success', 'Вы успешно вошли в систему');
      navigate('/dashboard');
    },
    onError: (err) => {
      addToast('error', getErrorMessage(err));
    },
  });
}

export function useRegister() {
  const { addToast } = useToastStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: (data: RegisterRequest) => authApi.register(data),
    onSuccess: () => {
      addToast('success', 'Регистрация успешна! Войдите в систему.');
      navigate('/login');
    },
    onError: (err) => {
      addToast('error', getErrorMessage(err));
    },
  });
}

export function useFetchUser() {
  const { accessToken, setUser } = useAuthStore();

  return useQuery({
    queryKey: ['auth', 'me'],
    queryFn: async () => {
      const res = await authApi.me();
      setUser(res.data);
      return res.data;
    },
    enabled: !!accessToken,
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
}
