import { apiClient } from './client';
import type { ProfileData, LeaderboardData } from '@/types';

export const profileApi = {
  getProfile: () =>
    apiClient.get<ProfileData>('/api/v1/profile'),

  getLeaderboard: () =>
    apiClient.get<LeaderboardData>('/api/v1/leaderboard'),
};
