import { useQuery } from '@tanstack/react-query';
import { coursesApi } from '@/api/courses';

export function useLesson(id: string) {
  return useQuery({
    queryKey: ['lessons', id],
    queryFn: async () => {
      const res = await coursesApi.getLesson(id);
      return res.data;
    },
    enabled: !!id,
    staleTime: 60 * 1000,
  });
}
