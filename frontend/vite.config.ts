import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');

  // Where is the backend gateway?
  // Strip /api/v1 if someone accidentally put it in the env
  const rawTarget = env.VITE_API_BASE_URL || 'http://localhost:8080';
  const proxyTarget = rawTarget.replace(/\/api\/v1\/?$/, '').replace(/\/+$/, '') || 'http://localhost:8080';

  return {
    plugins: [react()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      port: 3000,
      proxy: {
        '/api': {
          target: proxyTarget,
          changeOrigin: true,
          // No rewrite needed — /api/v1/* goes to gateway as /api/v1/*
        },
      },
    },
  };
});
