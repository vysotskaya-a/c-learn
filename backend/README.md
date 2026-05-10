# C-Learn Backend

Микросервисный бэкенд для платформы изучения языка C.

## Архитектура

```
Frontend → API Gateway (:8080)
             ├── Auth Service       (:8001)  — регистрация, JWT, роли
             ├── LMS Service        (:8002)  — курсы, уроки, задачи, CMS
             ├── Code Runner        (:8003)  — компиляция C в Docker
             └── Gamification       (:8004)  — XP, достижения, leaderboard
```

**БД:** PostgreSQL 16, одна инстанция с тремя схемами (`auth`, `lms`, `gamification`).
**Стек:** Go 1.22, Gin, JWT, bcrypt, Docker SDK.


## API Endpoints

### Auth (`/api/v1/auth`)
| Метод | URL | Auth |
|-------|-----|------|
| POST | `/register` | — |
| POST | `/login` | — |
| POST | `/refresh` | — |
| GET | `/me` | Bearer |

### LMS (`/api/v1`)
| Метод | URL | Auth |
|-------|-----|------|
| GET | `/courses/tree` | Bearer |
| GET | `/lessons/:id` | Bearer |
| POST | `/tasks/:id/run` | Bearer |
| POST | `/tasks/:id/submit` | Bearer |
| GET | `/solutions` | Bearer |

### CMS (`/api/v1/admin`, admin only)
| Метод | URL |
|-------|-----|
| GET/POST | `/modules` |
| PUT/DELETE | `/modules/:id` |
| POST | `/lessons` |
| PUT/DELETE | `/lessons/:id` |
| POST | `/tasks` |
| PUT/DELETE | `/tasks/:id` |
| PUT | `/tasks/:id/test-cases` |

### Gamification (`/api/v1`)
| Метод | URL | Auth |
|-------|-----|------|
| GET | `/profile` | Bearer |
| GET | `/leaderboard` | Bearer |

## Структура проекта

```
c-learn/
├── cmd/                    # Entrypoints (main.go для каждого сервиса)
├── internal/               # Бизнес-логика: handler → service → repository
│   ├── auth/               # Auth Service
│   ├── lms/                # LMS Service + CMS + Submit orchestration
│   ├── runner/             # Code Runner (Docker SDK)
│   ├── gamification/       # Gamification (XP, achievements)
│   ├── gateway/            # API Gateway (reverse proxy + JWT)
│   └── models/             # Shared DTO models
├── pkg/                    # Shared: response, validator, middleware, httpclient, errs
├── migrations/             # SQL-миграции (auth, lms, gamification)
├── deploy/                 # Dockerfiles + docker-compose.yml
├── scripts/                # migrate.sh, seed.sh, backup.sh
├── .github/workflows/      # CI: lint + test + build per service
├── go.mod
└── Makefile
```
