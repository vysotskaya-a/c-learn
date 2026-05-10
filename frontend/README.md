# C-Learn Frontend

Веб-приложение для обучения языку C с элементами геймификации.

## Стек технологий

- **React 18** + **TypeScript** + **Vite**
- **React Router v6** — маршрутизация
- **TanStack Query (React Query)** — серверное состояние, кеширование
- **Zustand** — клиентское состояние (auth, theme)
- **Monaco Editor** — редактор кода с подсветкой C
- **react-markdown** + remark-gfm + rehype-highlight — рендеринг теории
- **Tailwind CSS** — стилизация
- **Axios** — HTTP-клиент с interceptors

## Требования

- Node.js >= 18
- npm >= 9

## Структура проекта

```
src/
├── api/           # HTTP-клиент и API endpoints
├── components/
│   ├── ui/        # Переиспользуемые UI-компоненты
│   └── layout/    # Layout, Sidebar, ErrorBoundary
├── hooks/         # React Query хуки
├── pages/         # Страницы приложения
│   └── admin/     # CMS-страницы (admin only)
├── routes/        # Роутинг и guards
├── store/         # Zustand stores
├── types/         # TypeScript типы
└── utils/         # Утилиты
```

## Экраны

| Маршрут | Описание | Доступ |
|---------|----------|--------|
| `/login` | Вход | Публичный |
| `/register` | Регистрация | Публичный |
| `/dashboard` | Дерево курса, прогресс | Авторизованный |
| `/lessons/:id` | Рабочее пространство урока (split-view) | Авторизованный |
| `/profile` | Профиль, XP, достижения | Авторизованный |
| `/leaderboard` | Рейтинг | Авторизованный |
| `/admin/modules` | CMS: Модули | Admin |
| `/admin/lessons` | CMS: Уроки | Admin |
| `/admin/tasks` | CMS: Задачи + тест-кейсы | Admin |

## API-интеграция

Frontend ходит **только** через API Gateway на `/api/v1/*`.
Список используемых endpoints:

### Auth
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `GET  /api/v1/auth/me`

### LMS
- `GET  /api/v1/courses/tree`
- `GET  /api/v1/lessons/:id`
- `POST /api/v1/tasks/:id/run`
- `POST /api/v1/tasks/:id/submit`

### Gamification
- `GET /api/v1/profile`
- `GET /api/v1/leaderboard`

### Admin CMS
- `GET/POST       /api/v1/admin/modules`
- `PUT/DELETE      /api/v1/admin/modules/:id`
- `POST            /api/v1/admin/lessons`
- `PUT/DELETE      /api/v1/admin/lessons/:id`
- `POST            /api/v1/admin/tasks`
- `PUT/DELETE      /api/v1/admin/tasks/:id`
- `PUT             /api/v1/admin/tasks/:id/test-cases`


## Темы

Поддержка светлой и тёмной темы через Tailwind `dark:` class strategy.
Переключение в меню пользователя (sidebar).
