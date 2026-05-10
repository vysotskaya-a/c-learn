/* ========== Auth ========== */

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
}

export interface RegisterResponse {
  id: string;
  email: string;
  username: string;
  role: 'student' | 'admin';
  created_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface TokenResponse {
  access_token: string;
  refresh_token: string;
  token_type: 'Bearer';
  expires_in: number;
}

export interface RefreshRequest {
  refresh_token: string;
}

export interface UserInfo {
  id: string;
  email: string;
  username: string;
  role: 'student' | 'admin';
  created_at: string;
}

/* ========== LMS ========== */

export type LessonStatus = 'not_started' | 'in_progress' | 'completed';

export interface CourseLesson {
  id: string;
  title: string;
  sort_order: number;
  status: LessonStatus;
}

export interface CourseModule {
  id: string;
  title: string;
  sort_order: number;
  lessons: CourseLesson[];
}

export interface CourseTree {
  modules: CourseModule[];
}

export type ModuleStatus = 'locked' | 'available' | 'completed';

export interface TaskSample {
  input: string;
  expected: string;
}

export type Difficulty = 'easy' | 'medium' | 'hard';

export interface LessonTask {
  id: string;
  title: string;
  description: string;
  difficulty: Difficulty;
  is_solved: boolean;
  samples: TaskSample[];
}

export interface LessonDetail {
  id: string;
  module_id: string;
  title: string;
  theory_md: string;
  status: LessonStatus;
  tasks: LessonTask[];
}

/* ========== Code Run/Submit ========== */

export interface RunRequest {
  source_code: string;
  stdin: string;
}

export interface RunResponse {
  stdout: string;
  stderr: string;
  exit_code: number;
  exec_time_ms: number;
}

export interface SubmitRequest {
  source_code: string;
}

export type Verdict =
  | 'ok'
  | 'compilation_error'
  | 'wrong_answer'
  | 'time_limit_exceeded'
  | 'runtime_error';

export interface Achievement {
  code: string;
  title: string;
  icon_url: string;
}

export interface SubmitResponse {
  solution_id: string;
  verdict: Verdict;
  tests_passed: number;
  tests_total: number;
  exec_time_ms?: number;
  compiler_output?: string;
  failed_test?: number;
  xp_awarded: number;
  achievements_unlocked: Achievement[];
}

/* ========== Gamification ========== */

export interface ProfileAchievement {
  code: string;
  title: string;
  description: string;
  icon_url: string;
  awarded_at: string;
}

export interface ProfileData {
  user_id: string;
  total_xp: number;
  level: number;
  solved_count: number;
  achievements: ProfileAchievement[];
}

export interface LeaderboardEntry {
  rank: number;
  user_id: string;
  username?: string;
  total_xp: number;
  solved_count: number;
}

export interface LeaderboardData {
  entries: LeaderboardEntry[];
  total: number;
  current_user_rank: number;
}

/* ========== Admin / CMS ========== */

export interface AdminModule {
  id: string;
  title: string;
  description: string;
  sort_order: number;
}

export interface CreateModuleRequest {
  title: string;
  description: string;
  sort_order: number;
}

export interface CreateLessonRequest {
  module_id: string;
  title: string;
  theory_md: string;
  sort_order: number;
}

export interface AdminLesson {
  id: string;
  module_id: string;
  title: string;
  theory_md: string;
  sort_order: number;
}

export interface TestCaseInput {
  id?: string;
  input: string;
  expected: string;
  is_sample: boolean;
}

export interface CreateTaskRequest {
  lesson_id: string;
  title: string;
  description: string;
  difficulty: Difficulty;
  sort_order: number;
  test_cases: TestCaseInput[];
}

export interface AdminTask {
  id: string;
  lesson_id: string;
  title: string;
  description: string;
  difficulty: Difficulty;
  test_cases?: TestCaseInput[];
}

export interface UpdateTestCasesRequest {
  test_cases: TestCaseInput[];
}

export interface AdminFullTestCase {
  id: string;
  input: string;
  expected: string;
  is_sample: boolean;
  sort_order: number;
}

export interface AdminFullTask {
  id: string;
  lesson_id: string;
  title: string;
  description: string;
  difficulty: Difficulty;
  sort_order: number;
  test_cases: AdminFullTestCase[];
}

export interface AdminFullLesson {
  id: string;
  module_id: string;
  title: string;
  theory_md: string;
  sort_order: number;
  tasks: AdminFullTask[];
}

export interface AdminFullModule {
  id: string;
  title: string;
  description: string;
  sort_order: number;
  lessons: AdminFullLesson[];
}

/* ========== Errors ========== */

export interface ApiError {
  error: string;
  message: string;
  details?: Record<string, string>;
}
