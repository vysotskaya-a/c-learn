import { useQuery } from '@tanstack/react-query';
import { profileApi } from '@/api/profile';

export function useProfile() {
  return useQuery({
    queryKey: ['profile'],
    queryFn: async () => {
      const res = await profileApi.getProfile();
      return res.data;
    },
    staleTime: 60 * 1000,
  });
}

export function useLeaderboard() {
  return useQuery({
    queryKey: ['leaderboard'],
    queryFn: async () => {
      const res = await profileApi.getLeaderboard();
      return res.data;
    },
    staleTime: 60 * 1000,
    retry: 1,
  });
}
