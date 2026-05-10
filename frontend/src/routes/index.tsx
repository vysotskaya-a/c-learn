import { lazy, Suspense } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { ProtectedRoute, AdminRoute } from './guards';
import { AppLayout } from '@/components/layout/AppLayout';

// Lazy-loaded pages
const LoginPage = lazy(() => import('@/pages/LoginPage'));
const RegisterPage = lazy(() => import('@/pages/RegisterPage'));
const DashboardPage = lazy(() => import('@/pages/DashboardPage'));
const LessonPage = lazy(() => import('@/pages/LessonPage'));
const ProfilePage = lazy(() => import('@/pages/ProfilePage'));
const LeaderboardPage = lazy(() => import('@/pages/LeaderboardPage'));
const AdminModulesPage = lazy(() => import('@/pages/admin/AdminModulesPage'));
const AdminModuleEditorPage = lazy(() => import('@/pages/admin/AdminModuleEditorPage'));
const NotFoundPage = lazy(() => import('@/pages/NotFoundPage'));
const ForbiddenPage = lazy(() => import('@/pages/ForbiddenPage'));

function PageLoader() {
  return (
    <div className="flex items-center justify-center h-64">
      <div className="w-8 h-8 border-3 border-brand-500 border-t-transparent rounded-full animate-spin" />
    </div>
  );
}

function LayoutWrapper({ children }: { children: React.ReactNode }) {
  return (
    <ProtectedRoute>
      <AppLayout>{children}</AppLayout>
    </ProtectedRoute>
  );
}

export function AppRouter() {
  return (
    <Suspense fallback={<PageLoader />}>
      <Routes>
        {/* Public */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />

        {/* Protected with layout */}
        <Route
          path="/dashboard"
          element={
            <LayoutWrapper>
              <DashboardPage />
            </LayoutWrapper>
          }
        />
        <Route
          path="/lessons/:id"
          element={
            <ProtectedRoute>
              <LessonPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/profile"
          element={
            <LayoutWrapper>
              <ProfilePage />
            </LayoutWrapper>
          }
        />
        <Route
          path="/leaderboard"
          element={
            <LayoutWrapper>
              <LeaderboardPage />
            </LayoutWrapper>
          }
        />

        {/* Admin */}
        <Route
          path="/admin/modules"
          element={
            <AdminRoute>
              <AppLayout>
                <AdminModulesPage />
              </AppLayout>
            </AdminRoute>
          }
        />
        <Route
          path="/admin/modules/:id/edit"
          element={
            <AdminRoute>
              <AppLayout>
                <AdminModuleEditorPage />
              </AppLayout>
            </AdminRoute>
          }
        />

        {/* Redirects & fallbacks */}
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/forbidden" element={<ForbiddenPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </Suspense>
  );
}
