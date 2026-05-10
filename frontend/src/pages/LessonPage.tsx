import { useState, useCallback, lazy, Suspense, useMemo, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import { useLesson } from '@/hooks/useLesson';
import { useCourseTree } from '@/hooks/useCourseTree';
import { tasksApi } from '@/api/tasks';
import { getErrorMessage } from '@/api/client';
import { useToastStore } from '@/store/toastStore';
import { useThemeStore } from '@/store/themeStore';
import { exportCodeAsFile, isCodeTooLarge } from '@/utils/exportCode';
import { LessonSkeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import { Button } from '@/components/ui/Button';
import { Modal } from '@/components/ui/Modal';
import type { LessonTask, SubmitResponse, Verdict } from '@/types';
import {
  Play,
  Send,
  RotateCcw,
  Copy,
  Download,
  ChevronLeft,
  ChevronRight,
  CheckCircle2,
  XCircle,
  AlertTriangle,
  Clock,
  Zap,
  Award,
  Terminal,
  ClipboardPaste,
} from 'lucide-react';

const MonacoEditor = lazy(() => import('@monaco-editor/react'));

const DEFAULT_CODE = `#include <stdio.h>

int main() {
    // Ваш код здесь
    
    return 0;
}
`;

function DifficultyBadge({ difficulty }: { difficulty: string }) {
  const styles = {
    easy: 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400',
    medium: 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400',
    hard: 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400',
  };
  const labels = { easy: 'Легко', medium: 'Средне', hard: 'Сложно' };
  return (
    <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${styles[difficulty as keyof typeof styles] || ''}`}>
      {labels[difficulty as keyof typeof labels] || difficulty}
    </span>
  );
}

function VerdictIcon({ verdict }: { verdict: Verdict }) {
  switch (verdict) {
    case 'ok': return <CheckCircle2 size={48} className="text-emerald-500" />;
    case 'compilation_error': return <XCircle size={48} className="text-red-500" />;
    case 'wrong_answer': return <XCircle size={48} className="text-orange-500" />;
    case 'time_limit_exceeded': return <Clock size={48} className="text-amber-500" />;
    case 'runtime_error': return <AlertTriangle size={48} className="text-red-500" />;
  }
}

function verdictTitle(v: Verdict): string {
  switch (v) {
    case 'ok': return 'Все тесты пройдены ✓';
    case 'compilation_error': return 'Ошибка компиляции';
    case 'wrong_answer': return 'Неправильный ответ';
    case 'time_limit_exceeded': return 'Превышено время выполнения';
    case 'runtime_error': return 'Ошибка выполнения';
  }
}

export default function LessonPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { addToast } = useToastStore();
  const { dark } = useThemeStore();
  const { data: lesson, isLoading, error } = useLesson(id || '');
  const { data: tree } = useCourseTree();

  const [activeTaskIdx, setActiveTaskIdx] = useState(0);
  const [codes, setCodes] = useState<Record<string, string>>({});
  const [stdin, setStdin] = useState('');
  const [output, setOutput] = useState<{ stdout?: string; stderr?: string; execTime?: number } | null>(null);
  const [running, setRunning] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [submitResult, setSubmitResult] = useState<SubmitResponse | null>(null);
  const [showResultModal, setShowResultModal] = useState(false);

  useEffect(() => {
    setActiveTaskIdx(0);
    setSubmitResult(null);
    setShowResultModal(false);
    setOutput(null);
  }, [id]);

  const tasks = lesson?.tasks || [];
  const activeTask: LessonTask | undefined = tasks[activeTaskIdx];

  const currentCode = activeTask ? (codes[activeTask.id] ?? DEFAULT_CODE) : DEFAULT_CODE;

  const setCode = useCallback(
    (value: string | undefined) => {
      if (activeTask && value !== undefined) {
        setCodes((prev) => ({ ...prev, [activeTask.id]: value }));
      }
    },
    [activeTask]
  );

  // Find next lesson for navigation
  const nextLessonId = useMemo(() => {
    if (!tree || !lesson) return null;
    const allLessons = tree.modules
      .flatMap((m) => m.lessons.map((l) => ({ ...l, moduleSort: m.sort_order })))
      .sort((a, b) => a.moduleSort - b.moduleSort || a.sort_order - b.sort_order);
    const currentIdx = allLessons.findIndex((l) => l.id === lesson.id);
    if (currentIdx >= 0 && currentIdx < allLessons.length - 1) {
      return allLessons[currentIdx + 1].id;
    }
    return null;
  }, [tree, lesson]);

  const handleRun = async () => {
    if (!activeTask) return;
    if (!currentCode.trim()) {
      addToast('warning', 'Введите код перед запуском');
      return;
    }
    if (isCodeTooLarge(currentCode)) {
      addToast('error', 'Размер кода превышает 50 КБ');
      return;
    }

    setRunning(true);
    setOutput(null);
    try {
      const res = await tasksApi.run(activeTask.id, {
        source_code: currentCode,
        stdin,
      });
      setOutput({
        stdout: res.data.stdout,
        stderr: res.data.stderr,
        execTime: res.data.exec_time_ms,
      });
    } catch (err) {
      addToast('error', getErrorMessage(err));
    } finally {
      setRunning(false);
    }
  };

  const handleSubmit = async () => {
    if (!activeTask) return;
    if (!currentCode.trim()) {
      addToast('warning', 'Введите код перед отправкой');
      return;
    }
    if (isCodeTooLarge(currentCode)) {
      addToast('error', 'Размер кода превышает 50 КБ');
      return;
    }

    setSubmitting(true);
    try {
      const res = await tasksApi.submit(activeTask.id, { source_code: currentCode });
      setSubmitResult(res.data);
      setShowResultModal(true);

      // Invalidate caches to refresh progress
      queryClient.invalidateQueries({ queryKey: ['courses', 'tree'] });
      queryClient.invalidateQueries({ queryKey: ['lessons', id] });
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    } catch (err) {
      addToast('error', getErrorMessage(err));
    } finally {
      setSubmitting(false);
    }
  };

  const handleInsertSample = () => {
    if (activeTask?.samples?.[0]) {
      setStdin(activeTask.samples[0].input);
    }
  };

  if (isLoading) {
    return (
      <div className="h-full">
        <LessonSkeleton />
      </div>
    );
  }

  if (error || !lesson) {
    return (
      <div className="h-full flex items-center justify-center p-6">
        <EmptyState title="Урок не найден" description="Не удалось загрузить урок." />
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      {/* Top bar */}
      <div className="flex items-center gap-3 px-4 py-2.5 border-b border-surface-100 dark:border-surface-800 bg-white dark:bg-surface-900 shrink-0">
        <Button variant="ghost" size="sm" onClick={() => navigate('/dashboard')}>
          <ChevronLeft size={16} />
          Назад
        </Button>
        <h1 className="text-sm font-semibold text-surface-900 dark:text-surface-100 truncate flex-1">
          {lesson.title}
        </h1>
        {tasks.length > 1 && (
          <div className="flex items-center gap-1">
            {tasks.map((t, i) => (
              <button
                key={t.id}
                onClick={() => setActiveTaskIdx(i)}
                className={`
                  w-8 h-8 rounded-lg text-xs font-medium transition-colors flex items-center justify-center
                  ${i === activeTaskIdx
                    ? 'bg-brand-600 text-white'
                    : t.is_solved
                      ? 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400'
                      : 'bg-surface-100 dark:bg-surface-800 text-surface-600 dark:text-surface-400 hover:bg-surface-200 dark:hover:bg-surface-700'
                  }
                `}
                title={t.title}
              >
                {i + 1}
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Split view */}
      <div className="flex-1 flex overflow-hidden">
        {/* Left: Theory + Task */}
        <div className="w-1/2 border-r border-surface-100 dark:border-surface-800 overflow-y-auto">
          <div className="p-6 space-y-6">
            {/* Theory */}
            <div className="prose-content text-surface-800 dark:text-surface-200">
              <ReactMarkdown remarkPlugins={[remarkGfm]} rehypePlugins={[rehypeHighlight]}>
                {lesson.theory_md}
              </ReactMarkdown>
            </div>

            {/* Task description */}
            {activeTask && (
              <div className="border-t border-surface-100 dark:border-surface-800 pt-6">
                <div className="flex items-center gap-3 mb-3">
                  <h3 className="text-lg font-bold text-surface-900 dark:text-surface-100">
                    {activeTask.title}
                  </h3>
                  <DifficultyBadge difficulty={activeTask.difficulty} />
                  {activeTask.is_solved && (
                    <CheckCircle2 size={18} className="text-emerald-500" />
                  )}
                </div>
                <div className="prose-content text-surface-700 dark:text-surface-300 text-sm">
                  <ReactMarkdown remarkPlugins={[remarkGfm]} rehypePlugins={[rehypeHighlight]}>
                    {activeTask.description}
                  </ReactMarkdown>
                </div>

                {/* Samples */}
                {activeTask.samples.length > 0 && (
                  <div className="mt-4 space-y-3">
                    <h4 className="text-sm font-semibold text-surface-700 dark:text-surface-300">
                      Примеры
                    </h4>
                    {activeTask.samples.map((s, i) => (
                      <div
                        key={i}
                        className="grid grid-cols-2 gap-3 p-3 bg-surface-50 dark:bg-surface-800 rounded-lg"
                      >
                        <div>
                          <p className="text-xs font-medium text-surface-500 mb-1">Ввод</p>
                          <pre className="text-xs font-mono whitespace-pre-wrap bg-white dark:bg-surface-900 p-2 rounded border border-surface-100 dark:border-surface-700">
                            {s.input || '(пусто)'}
                          </pre>
                        </div>
                        <div>
                          <p className="text-xs font-medium text-surface-500 mb-1">Ожидаемый вывод</p>
                          <pre className="text-xs font-mono whitespace-pre-wrap bg-white dark:bg-surface-900 p-2 rounded border border-surface-100 dark:border-surface-700">
                            {s.expected}
                          </pre>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
        </div>

        {/* Right: Editor + Console */}
        <div className="w-1/2 flex flex-col bg-white dark:bg-surface-900">
          {/* Editor toolbar */}
          <div className="flex items-center gap-1 px-3 py-2 border-b border-surface-100 dark:border-surface-800">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setCode(DEFAULT_CODE)}
              title="Сбросить код"
            >
              <RotateCcw size={14} />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigator.clipboard.writeText(currentCode)}
              title="Копировать"
            >
              <Copy size={14} />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => exportCodeAsFile(currentCode)}
              title="Скачать .c"
            >
              <Download size={14} />
            </Button>
            <div className="flex-1" />
            <span className="text-xs text-surface-400 font-mono">main.c</span>
          </div>

          {/* Monaco Editor */}
          <div className="flex-1 min-h-0">
            <Suspense
              fallback={
                <div className="flex items-center justify-center h-full text-surface-400 text-sm">
                  Загрузка редактора...
                </div>
              }
            >
              <MonacoEditor
                height="100%"
                language="c"
                theme={dark ? 'vs-dark' : 'light'}
                value={currentCode}
                onChange={setCode}
                options={{
                  fontSize: 14,
                  fontFamily: '"JetBrains Mono", "Fira Code", monospace',
                  minimap: { enabled: false },
                  lineNumbers: 'on',
                  scrollBeyondLastLine: false,
                  wordWrap: 'on',
                  padding: { top: 12 },
                  tabSize: 4,
                  automaticLayout: true,
                }}
              />
            </Suspense>
          </div>

          {/* Stdin area */}
          <div className="border-t border-surface-100 dark:border-surface-800">
            <div className="flex items-center gap-2 px-3 py-1.5">
              <Terminal size={14} className="text-surface-400" />
              <span className="text-xs font-medium text-surface-500">stdin</span>
              {activeTask?.samples?.[0] && (
                <button
                  onClick={handleInsertSample}
                  className="ml-auto text-xs text-brand-600 dark:text-brand-400 flex items-center gap-1 hover:underline"
                >
                  <ClipboardPaste size={12} />
                  Вставить пример
                </button>
              )}
            </div>
            <textarea
              value={stdin}
              onChange={(e) => setStdin(e.target.value)}
              className="w-full px-3 pb-2 text-xs font-mono bg-transparent resize-none h-16 focus:outline-none text-surface-800 dark:text-surface-200 placeholder:text-surface-400"
              placeholder="Входные данные для запуска..."
            />
          </div>

          {/* Action buttons */}
          <div className="flex items-center gap-2 px-3 py-2.5 border-t border-surface-100 dark:border-surface-800">
            <Button
              variant="secondary"
              size="sm"
              onClick={handleRun}
              loading={running}
              disabled={submitting}
            >
              <Play size={14} />
              Запустить
            </Button>
            <Button
              variant="primary"
              size="sm"
              onClick={handleSubmit}
              loading={submitting}
              disabled={running}
            >
              <Send size={14} />
              Отправить
            </Button>
          </div>

          {/* Console output */}
          {output && (
            <div className="border-t border-surface-100 dark:border-surface-800 bg-surface-950 text-surface-200 p-3 max-h-48 overflow-y-auto">
              <div className="flex items-center gap-2 mb-2">
                <Terminal size={14} className="text-surface-400" />
                <span className="text-xs font-medium text-surface-400">Вывод</span>
                {output.execTime !== undefined && (
                  <span className="text-xs text-surface-500 ml-auto">{output.execTime} мс</span>
                )}
              </div>
              {output.stdout && (
                <pre className="text-xs font-mono whitespace-pre-wrap text-emerald-400 mb-1">
                  {output.stdout}
                </pre>
              )}
              {output.stderr && (
                <pre className="text-xs font-mono whitespace-pre-wrap text-red-400">
                  {output.stderr}
                </pre>
              )}
              {!output.stdout && !output.stderr && (
                <p className="text-xs text-surface-500 italic">Нет вывода</p>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Submit result modal */}
      <Modal
        open={showResultModal}
        onClose={() => setShowResultModal(false)}
        title="Результат проверки"
        maxWidth="max-w-md"
      >
        {submitResult && (
          <div className="text-center space-y-4">
            <VerdictIcon verdict={submitResult.verdict} />
            <h3
              className={`text-xl font-bold ${
                submitResult.verdict === 'ok' ? 'text-emerald-600 dark:text-emerald-400' : 'text-surface-900 dark:text-surface-100'
              }`}
            >
              {verdictTitle(submitResult.verdict)}
            </h3>

            <p className="text-sm text-surface-600 dark:text-surface-400">
              Пройдено {submitResult.tests_passed} из {submitResult.tests_total} тестов
            </p>

            {submitResult.verdict === 'wrong_answer' && submitResult.failed_test !== undefined && (
              <p className="text-sm text-orange-600 dark:text-orange-400">
                Первый проваленный тест: #{submitResult.failed_test}
              </p>
            )}

            {submitResult.verdict === 'compilation_error' && submitResult.compiler_output && (
              <pre className="text-left text-xs font-mono bg-surface-100 dark:bg-surface-800 p-3 rounded-lg overflow-x-auto text-red-600 dark:text-red-400 max-h-32 overflow-y-auto">
                {submitResult.compiler_output}
              </pre>
            )}

            {submitResult.verdict === 'ok' && (
              <div className="space-y-3 pt-2">
                {submitResult.xp_awarded > 0 && (
                  <div className="flex items-center justify-center gap-2 text-amber-600 dark:text-amber-400">
                    <Zap size={18} />
                    <span className="font-semibold">+{submitResult.xp_awarded} XP</span>
                  </div>
                )}
                {submitResult.achievements_unlocked.length > 0 && (
                  <div className="space-y-2">
                    {submitResult.achievements_unlocked.map((a) => (
                      <div
                        key={a.code}
                        className="flex items-center gap-2 justify-center bg-amber-50 dark:bg-amber-900/20 px-3 py-2 rounded-lg"
                      >
                        <Award size={16} className="text-amber-500" />
                        <span className="text-sm font-medium text-amber-800 dark:text-amber-300">
                          {a.title}
                        </span>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}

            <div className="flex justify-center gap-3 pt-4">
              <Button variant="secondary" onClick={() => setShowResultModal(false)}>
                Закрыть
              </Button>
              {submitResult.verdict === 'ok' && activeTaskIdx < tasks.length - 1 && (
                <Button
                  onClick={() => {
                    setShowResultModal(false);
                    setSubmitResult(null);
                    setActiveTaskIdx(activeTaskIdx + 1);
                  }}
                >
                  Следующая задача
                  <ChevronRight size={16} />
                </Button>
              )}
              {submitResult.verdict === 'ok' && activeTaskIdx >= tasks.length - 1 && nextLessonId && (
                <Button onClick={() => { setShowResultModal(false); navigate(`/lessons/${nextLessonId}`); }}>
                  Следующий урок
                  <ChevronRight size={16} />
                </Button>
              )}
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}
