import { useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  ChevronDown,
  ChevronRight,
  Plus,
  Save,
  Trash2,
  X,
  ArrowLeft,
  FileText,
  Code,
} from 'lucide-react';
import { adminApi } from '@/api/admin';
import { getErrorMessage } from '@/api/client';
import { useToastStore } from '@/store/toastStore';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import { ConfirmDialog } from '@/components/ui/ConfirmDialog';
import { CardSkeleton } from '@/components/ui/Skeleton';
import type {
  AdminFullLesson,
  AdminFullModule,
  AdminFullTask,
  AdminFullTestCase,
  Difficulty,
  CreateLessonRequest,
  CreateTaskRequest,
  CreateModuleRequest,
  TestCaseInput,
} from '@/types';

type LocalTestCase = AdminFullTestCase & { _local?: boolean };
type LocalTask = AdminFullTask & { _local?: boolean; _expanded?: boolean; test_cases: LocalTestCase[] };
type LocalLesson = AdminFullLesson & { _local?: boolean; _expanded?: boolean; tasks: LocalTask[] };

const DIFFICULTY_OPTIONS = [
  { value: 'easy', label: 'Легко' },
  { value: 'medium', label: 'Средне' },
  { value: 'hard', label: 'Сложно' },
];

function tempId() {
  return 'tmp-' + Math.random().toString(36).slice(2, 10);
}

function emptyTestCase(sortOrder: number, isSample = false): LocalTestCase {
  return {
    id: tempId(),
    input: '',
    expected: '',
    is_sample: isSample,
    sort_order: sortOrder,
    _local: true,
  };
}

function emptyTask(lessonID: string, sortOrder: number): LocalTask {
  return {
    id: tempId(),
    lesson_id: lessonID,
    title: '',
    description: '',
    difficulty: 'easy',
    sort_order: sortOrder,
    test_cases: [emptyTestCase(0, true)],
    _local: true,
    _expanded: true,
  };
}

function emptyLesson(moduleID: string, sortOrder: number): LocalLesson {
  return {
    id: tempId(),
    module_id: moduleID,
    title: '',
    theory_md: '',
    sort_order: sortOrder,
    tasks: [],
    _local: true,
    _expanded: true,
  };
}

export default function AdminModuleEditorPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { addToast } = useToastStore();

  const { data, isLoading, error } = useQuery<AdminFullModule>({
    queryKey: ['admin', 'module-full', id],
    queryFn: async () => (await adminApi.getModuleFull(id!)).data,
    enabled: !!id,
  });

  const [moduleTitle, setModuleTitle] = useState('');
  const [moduleDescription, setModuleDescription] = useState('');
  const [moduleSortOrder, setModuleSortOrder] = useState(0);
  const [lessons, setLessons] = useState<LocalLesson[]>([]);
  const [deleteLesson, setDeleteLesson] = useState<LocalLesson | null>(null);
  const [deleteTask, setDeleteTask] = useState<{ lessonId: string; task: LocalTask } | null>(null);

  useEffect(() => {
    if (!data) return;
    setModuleTitle(data.title);
    setModuleDescription(data.description);
    setModuleSortOrder(data.sort_order);
    setLessons((prev) => {
      const prevLessonMap = new Map<string, LocalLesson>(prev.map((l) => [l.id, l]));
      return data.lessons.map((l) => {
        const existing = prevLessonMap.get(l.id);
        const prevTaskMap = new Map<string, LocalTask>(
          existing ? existing.tasks.map((t) => [t.id, t]) : []
        );
        return {
          ...l,
          _expanded: existing?._expanded ?? false,
          tasks: l.tasks.map((t) => {
            const prevTask = prevTaskMap.get(t.id);
            return {
              ...t,
              _expanded: prevTask?._expanded ?? false,
              test_cases: [...t.test_cases],
            };
          }),
        };
      });
    });
  }, [data]);

  const saveModuleMutation = useMutation({
    mutationFn: (payload: CreateModuleRequest) =>
      adminApi.updateModule(id!, payload),
    onSuccess: () => {
      addToast('success', 'Модуль сохранён');
      queryClient.invalidateQueries({ queryKey: ['admin', 'modules'] });
      queryClient.invalidateQueries({ queryKey: ['admin', 'module-full', id] });
    },
    onError: (e) => addToast('error', getErrorMessage(e)),
  });

  const updateLessonLocal = (lessonId: string, patch: Partial<LocalLesson>) => {
    setLessons((prev) => prev.map((l) => (l.id === lessonId ? { ...l, ...patch } : l)));
  };

  const updateTaskLocal = (lessonId: string, taskId: string, patch: Partial<LocalTask>) => {
    setLessons((prev) =>
      prev.map((l) =>
        l.id === lessonId
          ? { ...l, tasks: l.tasks.map((t) => (t.id === taskId ? { ...t, ...patch } : t)) }
          : l
      )
    );
  };

  const addLesson = () => {
    const sortOrder = lessons.length > 0 ? Math.max(...lessons.map((l) => l.sort_order)) + 1 : 0;
    const fresh = emptyLesson(id!, sortOrder);
    setLessons((prev) => [...prev, fresh]);
  };

  const addTask = (lessonId: string) => {
    setLessons((prev) =>
      prev.map((l) => {
        if (l.id !== lessonId) return l;
        const sortOrder =
          l.tasks.length > 0 ? Math.max(...l.tasks.map((t) => t.sort_order)) + 1 : 0;
        return { ...l, tasks: [...l.tasks, emptyTask(lessonId, sortOrder)] };
      })
    );
  };

  const saveLesson = async (lesson: LocalLesson) => {
    if (!lesson.title.trim()) {
      addToast('error', 'Название урока обязательно');
      return;
    }
    try {
      const payload: CreateLessonRequest = {
        module_id: id!,
        title: lesson.title.trim(),
        theory_md: lesson.theory_md,
        sort_order: lesson.sort_order,
      };
      if (lesson._local) {
        const res = await adminApi.createLesson(payload);
        updateLessonLocal(lesson.id, {
          id: res.data.id,
          _local: false,
        });
        addToast('success', 'Урок создан');
      } else {
        await adminApi.updateLesson(lesson.id, payload);
        addToast('success', 'Урок обновлён');
      }
      queryClient.invalidateQueries({ queryKey: ['admin', 'module-full', id] });
      queryClient.invalidateQueries({ queryKey: ['courses', 'tree'] });
    } catch (e) {
      addToast('error', getErrorMessage(e));
    }
  };

  const removeLesson = async (lesson: LocalLesson) => {
    if (lesson._local) {
      setLessons((prev) => prev.filter((l) => l.id !== lesson.id));
      setDeleteLesson(null);
      return;
    }
    try {
      await adminApi.deleteLesson(lesson.id);
      setLessons((prev) => prev.filter((l) => l.id !== lesson.id));
      addToast('success', 'Урок удалён');
      queryClient.invalidateQueries({ queryKey: ['courses', 'tree'] });
    } catch (e) {
      addToast('error', getErrorMessage(e));
    } finally {
      setDeleteLesson(null);
    }
  };

  const saveTask = async (lessonId: string, task: LocalTask) => {
    if (!task.title.trim()) {
      addToast('error', 'Название задачи обязательно');
      return;
    }
    if (task.test_cases.length === 0) {
      addToast('error', 'Минимум 1 тест-кейс');
      return;
    }
    if (task.test_cases.some((tc) => !tc.expected.trim())) {
      addToast('error', 'Поле Expected обязательно для всех тестов');
      return;
    }

    const lesson = lessons.find((l) => l.id === lessonId);
    if (!lesson || lesson._local) {
      addToast('error', 'Сначала сохраните урок');
      return;
    }

    const tcPayload: TestCaseInput[] = task.test_cases.map((tc) => ({
      input: tc.input,
      expected: tc.expected,
      is_sample: tc.is_sample,
    }));

    try {
      if (task._local) {
        const payload: CreateTaskRequest = {
          lesson_id: lesson.id,
          title: task.title.trim(),
          description: task.description,
          difficulty: task.difficulty,
          sort_order: task.sort_order,
          test_cases: tcPayload,
        };
        const res = await adminApi.createTask(payload);
        updateTaskLocal(lessonId, task.id, {
          id: res.data.id,
          _local: false,
          test_cases: task.test_cases.map((tc) => ({ ...tc, _local: false })),
        });
        addToast('success', 'Задача создана');
      } else {
        await adminApi.updateTask(task.id, {
          lesson_id: lesson.id,
          title: task.title.trim(),
          description: task.description,
          difficulty: task.difficulty,
          sort_order: task.sort_order,
          test_cases: tcPayload,
        });
        await adminApi.updateTestCases(task.id, { test_cases: tcPayload });
        updateTaskLocal(lessonId, task.id, {
          test_cases: task.test_cases.map((tc) => ({ ...tc, _local: false })),
        });
        addToast('success', 'Задача обновлена');
      }
      queryClient.invalidateQueries({ queryKey: ['admin', 'module-full', id] });
    } catch (e) {
      addToast('error', getErrorMessage(e));
    }
  };

  const removeTask = async (lessonId: string, task: LocalTask) => {
    if (task._local) {
      setLessons((prev) =>
        prev.map((l) =>
          l.id === lessonId ? { ...l, tasks: l.tasks.filter((t) => t.id !== task.id) } : l
        )
      );
      setDeleteTask(null);
      return;
    }
    try {
      await adminApi.deleteTask(task.id);
      setLessons((prev) =>
        prev.map((l) =>
          l.id === lessonId ? { ...l, tasks: l.tasks.filter((t) => t.id !== task.id) } : l
        )
      );
      addToast('success', 'Задача удалена');
    } catch (e) {
      addToast('error', getErrorMessage(e));
    } finally {
      setDeleteTask(null);
    }
  };

  const sortedLessons = useMemo(
    () => [...lessons].sort((a, b) => a.sort_order - b.sort_order),
    [lessons]
  );

  if (isLoading) {
    return (
      <div className="max-w-5xl mx-auto p-6 space-y-4">
        <CardSkeleton />
        <CardSkeleton />
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="max-w-5xl mx-auto p-6">
        <p className="text-danger-500">Не удалось загрузить модуль</p>
        <Button variant="secondary" onClick={() => navigate('/admin/modules')} className="mt-4">
          <ArrowLeft size={16} /> К списку модулей
        </Button>
      </div>
    );
  }

  return (
    <div className="max-w-5xl mx-auto p-6 space-y-6 animate-fade-in">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="sm" onClick={() => navigate('/admin/modules')}>
          <ArrowLeft size={16} /> К модулям
        </Button>
      </div>

      {/* Module form */}
      <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-700 p-6 space-y-4">
        <h1 className="text-xl font-bold text-surface-900 dark:text-surface-100">Модуль</h1>
        <Input
          label="Название"
          value={moduleTitle}
          onChange={(e) => setModuleTitle(e.target.value)}
        />
        <div className="space-y-1.5">
          <label className="block text-sm font-medium text-surface-700 dark:text-surface-300">
            Описание
          </label>
          <textarea
            value={moduleDescription}
            onChange={(e) => setModuleDescription(e.target.value)}
            className="w-full px-3.5 py-2.5 rounded-lg border text-sm bg-white dark:bg-surface-900 border-surface-200 dark:border-surface-700 text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-brand-500/40 resize-none h-20"
          />
        </div>
        <Input
          label="Порядок сортировки"
          type="number"
          value={String(moduleSortOrder)}
          onChange={(e) => setModuleSortOrder(Number(e.target.value))}
        />
        <div className="flex justify-end">
          <Button
            onClick={() =>
              saveModuleMutation.mutate({
                title: moduleTitle.trim(),
                description: moduleDescription.trim(),
                sort_order: moduleSortOrder,
              })
            }
            loading={saveModuleMutation.isPending}
          >
            <Save size={16} /> Сохранить модуль
          </Button>
        </div>
      </div>

      {/* Lessons */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-surface-900 dark:text-surface-100">Уроки</h2>
          <Button onClick={addLesson} size="sm">
            <Plus size={14} /> Добавить урок
          </Button>
        </div>

        {sortedLessons.length === 0 && (
          <div className="bg-white dark:bg-surface-900 rounded-xl border border-dashed border-surface-200 dark:border-surface-700 p-8 text-center">
            <FileText size={32} className="mx-auto text-surface-400 mb-2" />
            <p className="text-sm text-surface-500">В этом модуле ещё нет уроков</p>
          </div>
        )}

        {sortedLessons.map((lesson) => (
          <LessonCard
            key={lesson.id}
            lesson={lesson}
            onToggle={() => updateLessonLocal(lesson.id, { _expanded: !lesson._expanded })}
            onChange={(patch) => updateLessonLocal(lesson.id, patch)}
            onSave={() => saveLesson(lesson)}
            onDelete={() => setDeleteLesson(lesson)}
            onAddTask={() => addTask(lesson.id)}
            onTaskChange={(taskId, patch) => updateTaskLocal(lesson.id, taskId, patch)}
            onTaskSave={(task) => saveTask(lesson.id, task)}
            onTaskDelete={(task) => setDeleteTask({ lessonId: lesson.id, task })}
          />
        ))}
      </div>

      <ConfirmDialog
        open={!!deleteLesson}
        onClose={() => setDeleteLesson(null)}
        onConfirm={() => deleteLesson && removeLesson(deleteLesson)}
        title="Удалить урок"
        message={`Удалить урок «${deleteLesson?.title || 'без названия'}»? Все задачи внутри тоже будут удалены.`}
      />
      <ConfirmDialog
        open={!!deleteTask}
        onClose={() => setDeleteTask(null)}
        onConfirm={() => deleteTask && removeTask(deleteTask.lessonId, deleteTask.task)}
        title="Удалить задачу"
        message={`Удалить задачу «${deleteTask?.task.title || 'без названия'}»?`}
      />
    </div>
  );
}

interface LessonCardProps {
  lesson: LocalLesson;
  onToggle: () => void;
  onChange: (patch: Partial<LocalLesson>) => void;
  onSave: () => void;
  onDelete: () => void;
  onAddTask: () => void;
  onTaskChange: (taskId: string, patch: Partial<LocalTask>) => void;
  onTaskSave: (task: LocalTask) => void;
  onTaskDelete: (task: LocalTask) => void;
}

function LessonCard({
  lesson,
  onToggle,
  onChange,
  onSave,
  onDelete,
  onAddTask,
  onTaskChange,
  onTaskSave,
  onTaskDelete,
}: LessonCardProps) {
  const sortedTasks = useMemo(
    () => [...lesson.tasks].sort((a, b) => a.sort_order - b.sort_order),
    [lesson.tasks]
  );

  return (
    <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-700 overflow-hidden">
      <div className="flex items-center gap-2 px-4 py-3 border-b border-surface-100 dark:border-surface-800">
        <button onClick={onToggle} className="p-1 text-surface-500 hover:text-surface-900 dark:hover:text-surface-100">
          {lesson._expanded ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
        </button>
        <div className="flex-1 min-w-0">
          <span className="text-xs text-surface-400 mr-2">#{lesson.sort_order}</span>
          <span className="text-sm font-medium text-surface-900 dark:text-surface-100 truncate">
            {lesson.title || <em className="text-surface-400">Новый урок</em>}
          </span>
          {lesson._local && (
            <span className="ml-2 text-xs text-amber-600 dark:text-amber-400">не сохранён</span>
          )}
        </div>
        <span className="text-xs text-surface-500">
          {lesson.tasks.length} задач{lesson.tasks.length === 1 ? 'а' : lesson.tasks.length >= 2 && lesson.tasks.length <= 4 ? 'и' : ''}
        </span>
        <Button variant="ghost" size="sm" onClick={onDelete}>
          <Trash2 size={14} className="text-danger-500" />
        </Button>
      </div>

      {lesson._expanded && (
        <div className="p-4 space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div className="md:col-span-2">
              <Input
                label="Название урока"
                value={lesson.title}
                onChange={(e) => onChange({ title: e.target.value })}
              />
            </div>
            <Input
              label="Порядок"
              type="number"
              value={String(lesson.sort_order)}
              onChange={(e) => onChange({ sort_order: Number(e.target.value) })}
            />
          </div>
          <div className="space-y-1.5">
            <label className="block text-sm font-medium text-surface-700 dark:text-surface-300">
              Теория (Markdown)
            </label>
            <textarea
              value={lesson.theory_md}
              onChange={(e) => onChange({ theory_md: e.target.value })}
              placeholder="# Заголовок&#10;&#10;Текст теории..."
              className="w-full px-3 py-2 rounded-lg border text-sm font-mono bg-white dark:bg-surface-900 border-surface-200 dark:border-surface-700 text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-brand-500/40 resize-none h-40"
            />
          </div>
          <div className="flex justify-end">
            <Button onClick={onSave} size="sm">
              <Save size={14} /> {lesson._local ? 'Создать урок' : 'Сохранить урок'}
            </Button>
          </div>

          {/* Tasks */}
          <div className="border-t border-surface-100 dark:border-surface-800 pt-4 space-y-3">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-semibold text-surface-700 dark:text-surface-300">
                Задачи в уроке
              </h3>
              <Button size="sm" variant="secondary" onClick={onAddTask} disabled={lesson._local}>
                <Plus size={14} /> Добавить задачу
              </Button>
            </div>
            {lesson._local && (
              <p className="text-xs text-surface-500">Сначала сохраните урок, чтобы добавить задачи.</p>
            )}
            {sortedTasks.length === 0 && !lesson._local && (
              <div className="rounded-lg border border-dashed border-surface-200 dark:border-surface-700 p-4 text-center">
                <Code size={20} className="mx-auto text-surface-400 mb-1" />
                <p className="text-xs text-surface-500">Задач пока нет</p>
              </div>
            )}
            {sortedTasks.map((task) => (
              <TaskCard
                key={task.id}
                task={task}
                onToggle={() => onTaskChange(task.id, { _expanded: !task._expanded })}
                onChange={(patch) => onTaskChange(task.id, patch)}
                onSave={() => onTaskSave(task)}
                onDelete={() => onTaskDelete(task)}
              />
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

interface TaskCardProps {
  task: LocalTask;
  onToggle: () => void;
  onChange: (patch: Partial<LocalTask>) => void;
  onSave: () => void;
  onDelete: () => void;
}

function TaskCard({ task, onToggle, onChange, onSave, onDelete }: TaskCardProps) {
  const updateTC = (idx: number, patch: Partial<LocalTestCase>) => {
    onChange({
      test_cases: task.test_cases.map((tc, i) => (i === idx ? { ...tc, ...patch } : tc)),
    });
  };

  const addTC = () => {
    onChange({
      test_cases: [...task.test_cases, emptyTestCase(task.test_cases.length, false)],
    });
  };

  const removeTC = (idx: number) => {
    if (task.test_cases.length <= 1) return;
    onChange({ test_cases: task.test_cases.filter((_, i) => i !== idx) });
  };

  return (
    <div className="rounded-lg border border-surface-200 dark:border-surface-700 bg-surface-50/50 dark:bg-surface-800/40">
      <div className="flex items-center gap-2 px-3 py-2">
        <button onClick={onToggle} className="p-1 text-surface-500 hover:text-surface-900 dark:hover:text-surface-100">
          {task._expanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
        </button>
        <div className="flex-1 min-w-0">
          <span className="text-xs text-surface-400 mr-2">#{task.sort_order}</span>
          <span className="text-sm text-surface-900 dark:text-surface-100 truncate">
            {task.title || <em className="text-surface-400">Новая задача</em>}
          </span>
          {task._local && (
            <span className="ml-2 text-xs text-amber-600 dark:text-amber-400">не сохранена</span>
          )}
        </div>
        <span
          className={`px-2 py-0.5 rounded-full text-xs font-medium ${
            task.difficulty === 'easy'
              ? 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400'
              : task.difficulty === 'medium'
                ? 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400'
                : 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400'
          }`}
        >
          {task.difficulty === 'easy' ? 'Легко' : task.difficulty === 'medium' ? 'Средне' : 'Сложно'}
        </span>
        <Button variant="ghost" size="sm" onClick={onDelete}>
          <Trash2 size={12} className="text-danger-500" />
        </Button>
      </div>

      {task._expanded && (
        <div className="px-3 pb-3 space-y-3">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div className="md:col-span-2">
              <Input
                label="Название"
                value={task.title}
                onChange={(e) => onChange({ title: e.target.value })}
              />
            </div>
            <Input
              label="Порядок"
              type="number"
              value={String(task.sort_order)}
              onChange={(e) => onChange({ sort_order: Number(e.target.value) })}
            />
          </div>
          <div className="space-y-1.5">
            <label className="block text-sm font-medium text-surface-700 dark:text-surface-300">
              Описание
            </label>
            <textarea
              value={task.description}
              onChange={(e) => onChange({ description: e.target.value })}
              className="w-full px-3 py-2 rounded-lg border text-sm bg-white dark:bg-surface-900 border-surface-200 dark:border-surface-700 text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-brand-500/40 resize-none h-24"
              placeholder="Описание задачи (Markdown)..."
            />
          </div>
          <Select
            label="Сложность"
            options={DIFFICULTY_OPTIONS}
            value={task.difficulty}
            onChange={(e) => onChange({ difficulty: e.target.value as Difficulty })}
          />

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <label className="text-sm font-medium text-surface-700 dark:text-surface-300">
                Тест-кейсы
              </label>
              <Button variant="ghost" size="sm" onClick={addTC}>
                <Plus size={12} /> Добавить
              </Button>
            </div>
            {task.test_cases.map((tc, i) => (
              <div
                key={tc.id}
                className="p-2 border border-surface-200 dark:border-surface-700 rounded-lg space-y-2 bg-white dark:bg-surface-900"
              >
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-surface-500">Тест #{i + 1}</span>
                  <div className="flex items-center gap-2">
                    <label className="flex items-center gap-1.5 text-xs text-surface-600 dark:text-surface-400">
                      <input
                        type="checkbox"
                        checked={tc.is_sample}
                        onChange={(e) => updateTC(i, { is_sample: e.target.checked })}
                        className="rounded"
                      />
                      Пример
                    </label>
                    {task.test_cases.length > 1 && (
                      <button
                        onClick={() => removeTC(i)}
                        className="text-surface-400 hover:text-danger-500"
                      >
                        <X size={12} />
                      </button>
                    )}
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <label className="text-xs text-surface-500">Input</label>
                    <textarea
                      value={tc.input}
                      onChange={(e) => updateTC(i, { input: e.target.value })}
                      className="w-full px-2 py-1.5 rounded border text-xs font-mono bg-white dark:bg-surface-900 border-surface-200 dark:border-surface-700 text-surface-900 dark:text-surface-100 resize-none h-16"
                    />
                  </div>
                  <div>
                    <label className="text-xs text-surface-500">Expected *</label>
                    <textarea
                      value={tc.expected}
                      onChange={(e) => updateTC(i, { expected: e.target.value })}
                      className="w-full px-2 py-1.5 rounded border text-xs font-mono bg-white dark:bg-surface-900 border-surface-200 dark:border-surface-700 text-surface-900 dark:text-surface-100 resize-none h-16"
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>

          <div className="flex justify-end">
            <Button onClick={onSave} size="sm">
              <Save size={14} /> {task._local ? 'Создать задачу' : 'Сохранить задачу'}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
