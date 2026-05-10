import { apiClient } from './client';
import type {
  AdminModule,
  CreateModuleRequest,
  AdminLesson,
  CreateLessonRequest,
  AdminTask,
  CreateTaskRequest,
  UpdateTestCasesRequest,
  AdminFullModule,
} from '@/types';

export const adminApi = {
  // Modules
  getModules: () =>
    apiClient.get<AdminModule[]>('/api/v1/admin/modules'),

  getModuleFull: (id: string) =>
    apiClient.get<AdminFullModule>(`/api/v1/admin/modules/${id}`),

  createModule: (data: CreateModuleRequest) =>
    apiClient.post<AdminModule>('/api/v1/admin/modules', data),

  updateModule: (id: string, data: Partial<CreateModuleRequest>) =>
    apiClient.put<AdminModule>(`/api/v1/admin/modules/${id}`, data),

  deleteModule: (id: string) =>
    apiClient.delete(`/api/v1/admin/modules/${id}`),

  // Lessons
  createLesson: (data: CreateLessonRequest) =>
    apiClient.post<AdminLesson>('/api/v1/admin/lessons', data),

  updateLesson: (id: string, data: Partial<CreateLessonRequest>) =>
    apiClient.put<AdminLesson>(`/api/v1/admin/lessons/${id}`, data),

  deleteLesson: (id: string) =>
    apiClient.delete(`/api/v1/admin/lessons/${id}`),

  // Tasks
  createTask: (data: CreateTaskRequest) =>
    apiClient.post<AdminTask>('/api/v1/admin/tasks', data),

  updateTask: (id: string, data: Partial<CreateTaskRequest>) =>
    apiClient.put<AdminTask>(`/api/v1/admin/tasks/${id}`, data),

  deleteTask: (id: string) =>
    apiClient.delete(`/api/v1/admin/tasks/${id}`),

  updateTestCases: (taskId: string, data: UpdateTestCasesRequest) =>
    apiClient.put(`/api/v1/admin/tasks/${taskId}/test-cases`, data),
};
