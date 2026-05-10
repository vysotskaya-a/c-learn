import { apiClient } from './client';
import type { CourseTree, LessonDetail } from '@/types';

export const coursesApi = {
  getTree: () =>
    apiClient.get<CourseTree>('/api/v1/courses/tree'),

  getLesson: (id: string) =>
    apiClient.get<LessonDetail>(`/api/v1/lessons/${id}`),
};
