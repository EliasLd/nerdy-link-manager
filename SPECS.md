# Personal Link Manager Specification

## Objective

Develop a Go backend for a personal link manager that is secure, lighweight and extensible.
It will expose a **REST API** consumed by a **Svelte** + **TypeScript** frontend, using a SQLite database to store all data.

The application is designed for:
 - Single-user usage
 - Deployment & self-hosting
 - Providing link management and statistics

---

## Technical Stack

|  Component        | Technology        |
|-------------------|-------------------|
| Backend           | Go (no framework) |  
| Database          | SQLite            |
| Authentication    | JWT               |
| Migrations        | golang-migrate    |
| Containerization  | Docker/Podman     |
| Frontend          | Svelte + TS       |

## Backend Specification

### Architecture

```
backend/
├── cmd/
│ └── main.go
├── config/
│ └── config.go
├── internal/
│ ├── crypto/
│ │ ├── password.go
│ │ └── jwt.go
│ ├── db/
│ │ ├── db.go
│ │ ├── migrations/
│ │ └── models.go
│ ├── repositories/
│ │ ├── user_repository.go
│ │ ├── link_repository.go
│ │ └── stats_repository.go
│ ├── services/
│ │ ├── auth_service.go
│ │ ├── link_service.go
│ │ └── stats_service.go
│ ├── handlers/
│ │ ├── auth_handler.go
│ │ ├── link_handler.go
│ │ └── stats_handler.go
│ ├── middleware/
│ │ └── auth_middleware.go
│ └── router/
│   └── router.go
└── Dockerfile
```

### Authentication

The objective is to protect all API routes so that only the authorized user can:
 - Log in
 - Manage links
 - View statistics

Method
 - Email + password login
 - Password hashed with `bcrypt`
 - JWT generated after login
 - JWT snet via HTTP-only cookie or `Authorization` header

### API

#### Auth


`POST /api/auth/login`

**Body**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response**
```json
{
  "token": "jwt.token.here"
}
```

#### Links

`GET /api/links`
Returns all links

**Response**
```json
[
  {
    "id": 1,
    "title": "Google",
    "url": "https://google.com",
    "description": "Search engine",
    "creation_at": "2026-01-01T12:00:00Z",
    "updated_at": "2026-01-01T12:00:00Z"
  }
]
```

`POST /api/links`
Adds a new link.

**Body**
```json
{
  "title": "Github",
  "url":  "https://github.com",
  "description":  "Code hosting platform"
}
```

`PUT /api/links/{id}`
Updates a link

**Body**
```json
{
  "title": "Github",
  "url": "https://github.com",
  "description": "Developer platform"
}
```

`DELETE /api/links/{id}`
Deletes a link.

**Redirection + Tracking**

`GET /r/{id}`
 - Increments the clicks statistics
 - Redirects to the actual URL

**Statistics**

`GET /r/{id}`
Returns click statistics for a link

**Response example**

```json
{
  "total_clicks": 42,
  "clicks_by_day": [
    { "date": "2026-01-10", "count": 5 },
    { "date": "2026-01-11", "count": 7 }
  ]
}
```

### Business Logic

#### Services

`AuthService`
 - `Login(email, password) (token, error)`
 - `ValidateToken(token) (user, error)`

`LinkService`
 - `CreateLink(input)`
 - `GetAllLinks()`
 - `GetLinkByID(id)`
 - `UpdateLink(id, input)`
 - `DeleteLink(id)`

`StatsService`
 - `TrackClick(linkID)`
 - `GetStats(linkID)`
 - `GetStatsByPeriod(linkID, form, to)` (optional)

### Middleware

`AuthMiddleware`
 - Checks for JWT presence
 - Validates the token
 - Blocks access to protected routes if invalid

## Database

### Table `users`

There will be three main tables

- `users`: only one user is expected (`id`, `email`, `password_hash`, `created_at`)
- `links`: (`id`, `title`, `url`, `description`, `created_at`, `updated_at`)
- `link_clicks`: (`id`, `link_id`, `clicked_at`)

## Frontend

| Component         | Technology            |
|-------------------|-----------------------|
| Framework         | SvelteKit             |
| Language          | TypeScript            |
| Styling           | Tailwind CSS          |
| State Management  | Svelte stores         |
| API Communication | Fetch API             |
| Authentication    | JWT                   |
| Build Tool        | Vite (via SvelteKit)  |
| Containerization  | Docker/Podman         |


---

### Architecture

```text
frontend/
├── src/
│ ├── routes/
│ │ ├── +layout.svelte
│ │ ├── +layout.ts
│ │ ├── +page.svelte # Landing / Login
│ │ ├── dashboard/
│ │ │ ├── +page.svelte # Link list view
│ │ │ ├── +page.ts
│ │ │ ├── stats/
│ │ │ └── +page.svelte # Statistics page
│ ├── lib/
│ │ ├── components/
│ │ │ ├── atoms/
│ │ │ │ ├── Button.svelte
│ │ │ │ ├── Input.svelte
│ │ │ │ ├── Label.svelte
│ │ │ │ ├── Card.svelte
│ │ │ │ └── Divider.svelte
│ │ │ ├── molecules/
│ │ │ │ ├── LinkItem.svelte
│ │ │ │ ├── LinkForm.svelte
│ │ │ │ ├── StatCard.svelte
│ │ │ │ └── Navbar.svelte
│ │ │ ├── organisms/
│ │ │ │ ├── LinkList.svelte
│ │ │ │ ├── StatsPanel.svelte
│ │ │ │ └── AuthForm.svelte
│ │ │ └── layouts/
│ │ │ ├── DashboardLayout.svelte
│ │ │ └── AuthLayout.svelte
│ │ ├── stores/
│ │ │ ├── auth.store.ts
│ │ │ ├── links.store.ts
│ │ │ └── stats.store.ts
│ │ ├── api/
│ │ │ ├── client.ts
│ │ │ ├── auth.api.ts
│ │ │ ├── links.api.ts
│ │ │ └── stats.api.ts
│ │ ├── types/
│ │ │ ├── auth.types.ts
│ │ │ ├── link.types.ts
│ │ │ └── stats.types.ts
│ │ └── utils/
│ │ ├── validators.ts
│ │ └── formatters.ts
│ └── logo.txt
└── Dockerfile
```

### Atomic Design System

#### Atoms
Basic UI building blocks:
- `Button`
- `Input`
- `Label`
- `Card`
- `Divider`

No business logic — only styling and accessibility.

---

#### Molecules
Combinations of atoms with minimal logic:
- `LinkItem` (displays a single link)
- `LinkForm` (form to create/edit links)
- `StatCard` (displays a single stat metric)
- `Navbar`

---

#### Organisms
Complex UI blocks composed of molecules:
- `LinkList` (list of all links)
- `StatsPanel` (dashboard of stats)
- `AuthForm` (login form)

---

####  Layouts
Reusable page layouts:
- `AuthLayout` (centered login UI)
- `DashboardLayout` (navbar + content area)

---

### Pages & Routing

| Route               | Description               |
|---------------------|---------------------------|
| `/`                 | Login page                |
| `/dashboard`        | Main link management page |
| `/dashboard/stats`  | Statistics page           |

Routing handled by SvelteKit’s file-based routing.

---

### Authentication Flow

1. User lands on `/`
2. User submits login form
3. Frontend sends credentials to `POST /api/auth/login`
4. Backend returns JWT (cookie or token)
5. Frontend updates `auth.store`
6. User is redirected to `/dashboard`
7. Protected routes check auth store before rendering

---

### State Management (Stores)

#### `auth.store.ts`
- Holds authentication state (`isAuthenticated`, `user`)
- Exposes:
  - `login()`
  - `checkAuth()`

---

#### `links.store.ts`
- Holds list of links
- Exposes:
  - `fetchLinks()`
  - `createLink(data)`
  - `updateLink(id, data)`
  - `deleteLink(id)`

---

#### `stats.store.ts`
- Holds statistics data
- Exposes:
  - `fetchStats(linkId)`
  - `fetchGlobalStats()` (optional)

---

### API Layer

#### `api/client.ts`
Centralized fetch wrapper:
- Automatically attaches auth token
- Handles errors and response parsing

---

#### API Modules

- `auth.api.ts`
  - `login(credentials)`
  - `logout()`

- `links.api.ts`
  - `getLinks()`
  - `createLink(data)`
  - `updateLink(id, data)`
  - `deleteLink(id)`

- `stats.api.ts`
  - `getStats(linkId)`
