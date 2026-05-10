import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminApi } from '@/api/admin';
import { getErrorMessage } from '@/api/client';
import { useToastStore } from '@/store/toastStore';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';
import { ConfirmDialog } from '@/components/ui/ConfirmDialog';
import { CardSkeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import { Plus, Pencil, Trash2, Package, FileText } from 'lucide-react';
import type { AdminModule, CreateModuleRequest } from '@/types';

export default function AdminModulesPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { addToast } = useToastStore();

  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<AdminModule | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<AdminModule | null>(null);

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [sortOrder, setSortOrder] = useState(0);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const { data: modules, isLoading } = useQuery({
    queryKey: ['admin', 'modules'],
    queryFn: async () => {
      const res = await adminApi.getModules();
      return res.data;
    },
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateModuleRequest) => adminApi.createModule(data),
    onSuccess: (res) => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'modules'] });
      addToast('success', 'Модуль создан');
      closeModal();
      navigate(`/admin/modules/${res.data.id}/edit`);
    },
    onError: (e) => addToast('error', getErrorMessage(e)),
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateModuleRequest> }) =>
      adminApi.updateModule(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'modules'] });
      addToast('success', 'Модуль обновлён');
      closeModal();
    },
    onError: (e) => addToast('error', getErrorMessage(e)),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => adminApi.deleteModule(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'modules'] });
      addToast('success', 'Модуль удалён');
      setDeleteTarget(null);
    },
    onError: (e) => addToast('error', getErrorMessage(e)),
  });

  const openCreate = () => {
    setEditing(null);
    setTitle('');
    setDescription('');
    setSortOrder(0);
    setErrors({});
    setModalOpen(true);
  };

  const openEdit = (mod: AdminModule) => {
    setEditing(mod);
    setTitle(mod.title);
    setDescription(mod.description || '');
    setSortOrder(mod.sort_order);
    setErrors({});
    setModalOpen(true);
  };

  const closeModal = () => {
    setModalOpen(false);
    setEditing(null);
  };

  const validate = (): boolean => {
    const e: Record<string, string> = {};
    if (!title.trim()) e.title = 'Название обязательно';
    if (sortOrder < 0) e.sortOrder = 'Порядок >= 0';
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleSubmit = () => {
    if (!validate()) return;
    const data: CreateModuleRequest = {
      title: title.trim(),
      description: description.trim(),
      sort_order: sortOrder,
    };
    if (editing) {
      updateMutation.mutate({ id: editing.id, data });
    } else {
      createMutation.mutate(data);
    }
  };

  if (isLoading) {
    return (
      <div className="max-w-4xl mx-auto p-6 space-y-4">
        <CardSkeleton />
        <CardSkeleton />
      </div>
    );
  }

  const sortedModules = [...(modules || [])].sort((a, b) => a.sort_order - b.sort_order);

  return (
    <div className="max-w-4xl mx-auto p-6 space-y-6 animate-fade-in">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-surface-900 dark:text-surface-100">Модули</h1>
        <Button onClick={openCreate}>
          <Plus size={16} /> Создать модуль
        </Button>
      </div>

      {sortedModules.length === 0 ? (
        <EmptyState
          icon={<Package size={40} />}
          title="Нет модулей"
          description="Создайте первый модуль курса."
          action={
            <Button onClick={openCreate}>
              <Plus size={16} /> Создать
            </Button>
          }
        />
      ) : (
        <div className="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-700 overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="border-b border-surface-100 dark:border-surface-800">
                <th className="text-left px-5 py-3 text-xs font-semibold text-surface-500 uppercase w-16">#</th>
                <th className="text-left px-5 py-3 text-xs font-semibold text-surface-500 uppercase">Название</th>
                <th className="text-left px-5 py-3 text-xs font-semibold text-surface-500 uppercase hidden md:table-cell">Описание</th>
                <th className="text-right px-5 py-3 text-xs font-semibold text-surface-500 uppercase w-32">Действия</th>
              </tr>
            </thead>
            <tbody>
              {sortedModules.map((mod) => (
                <tr key={mod.id} className="border-b border-surface-50 dark:border-surface-800/50 last:border-0">
                  <td className="px-5 py-3 text-sm text-surface-500">{mod.sort_order}</td>
                  <td className="px-5 py-3 text-sm font-medium text-surface-900 dark:text-surface-100">
                    <button
                      onClick={() => navigate(`/admin/modules/${mod.id}/edit`)}
                      className="text-left hover:text-brand-600 dark:hover:text-brand-400 transition-colors"
                    >
                      {mod.title}
                    </button>
                  </td>
                  <td className="px-5 py-3 text-sm text-surface-500 hidden md:table-cell truncate max-w-[200px]">
                    {mod.description || '—'}
                  </td>
                  <td className="px-5 py-3 text-right">
                    <div className="flex justify-end gap-1">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => navigate(`/admin/modules/${mod.id}/edit`)}
                        title="Уроки и задачи"
                      >
                        <FileText size={14} />
                      </Button>
                      <Button variant="ghost" size="sm" onClick={() => openEdit(mod)} title="Редактировать">
                        <Pencil size={14} />
                      </Button>
                      <Button variant="ghost" size="sm" onClick={() => setDeleteTarget(mod)} title="Удалить">
                        <Trash2 size={14} className="text-danger-500" />
                      </Button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Create/Edit Modal */}
      <Modal open={modalOpen} onClose={closeModal} title={editing ? 'Редактировать модуль' : 'Создать модуль'}>
        <div className="space-y-4">
          <Input label="Название" value={title} onChange={(e) => setTitle(e.target.value)} error={errors.title} />
          <div className="space-y-1.5">
            <label className="block text-sm font-medium text-surface-700 dark:text-surface-300">Описание</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full px-3.5 py-2.5 rounded-lg border text-sm bg-white dark:bg-surface-900 border-surface-200 dark:border-surface-700 text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 resize-none h-24"
            />
          </div>
          <Input
            label="Порядок сортировки"
            type="number"
            value={String(sortOrder)}
            onChange={(e) => setSortOrder(Number(e.target.value))}
            error={errors.sortOrder}
          />
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="secondary" onClick={closeModal}>
              Отмена
            </Button>
            <Button
              onClick={handleSubmit}
              loading={createMutation.isPending || updateMutation.isPending}
            >
              {editing ? 'Сохранить' : 'Создать'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* Delete confirm */}
      <ConfirmDialog
        open={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        onConfirm={() => deleteTarget && deleteMutation.mutate(deleteTarget.id)}
        title="Удалить модуль"
        message={`Вы уверены, что хотите удалить модуль "${deleteTarget?.title}"? Это действие необратимо.`}
        loading={deleteMutation.isPending}
      />
    </div>
  );
}
