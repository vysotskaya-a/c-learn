import { useQuery } from '@tanstack/react-query';
import { coursesApi } from '@/api/courses';

export function useCourseTree() {
  return useQuery({
    queryKey: ['courses', 'tree'],
    queryFn: async () => {
      const res = await coursesApi.getTree();
      return res.data;
    },
    staleTime: 2 * 60 * 1000,
  });
}
