import { apiClient } from './client';
import type { RunRequest, RunResponse, SubmitRequest, SubmitResponse } from '@/types';

export const tasksApi = {
  run: (taskId: string, data: RunRequest) =>
    apiClient.post<RunResponse>(`/api/v1/tasks/${taskId}/run`, data),

  submit: (taskId: string, data: SubmitRequest) =>
    apiClient.post<SubmitResponse>(`/api/v1/tasks/${taskId}/submit`, data),
};
