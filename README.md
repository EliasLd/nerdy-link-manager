# Nerdy Link Manager

A self-hostable, lightweight, and flexible personal link manager built with Go and Svelte.

## Features

- Link management: create, edit, delete, and open links
- Folder organization: create/rename/delete folders and move links between folders
- Drag and drop: quickly move links into folders (and back to the main dashboard)
- Fast search: fuzzy search across all links with quick open
- Custom icons: automatic favicon detection with optional custom icon upload
- Authentication: JWT-based auth and protected API endpoints
- Clean dashboard UX: responsive layout and keyboard-friendly interactions
- Lightweight and simple: minimal dependencies, easy to self-host

## Installation / Configuration

The project ships with a guided installer script that builds and runs the full stack (backend + frontend) using Podman.

### Quick start

From the repository root:

```bash
chmod +x run.sh
./run.sh
```

The script will:
1. Request sudo rights (needed for host setup depending on your environment).
2. Ensure Podman is installed.
3. Load configuration from a `.env` file located at the repository root (`./.env`).
   - If `./.env` does not exist, the script will offer to create it interactively.
4. Detect a first run (database not found) and require initial admin credentials.
5. Build the backend and frontend container images.
6. Create a Podman pod (ports + persistent volume) and start both containers.

### Environment variables

The installer uses `./.env` (repository root). You can either:
- let `install.sh` create it interactively, or
- create it manually before running the script.

Required:
- `PUBLIC_API_URL` (used at frontend build time, e.g. `https://lm.example.com`)
- `JWT_SECRET` (backend JWT signing secret)

Required on first run only (when the database does not exist yet):
- `INITIAL_ADMIN_EMAIL`
- `INITIAL_ADMIN_PASSWORD`

Optional (with defaults):
- `FRONTEND_PORT` (default: `20301`)
- `BACKEND_PORT` (default: `9001`)
- `DATA_DIR` (default: `./backend/data`, mounted as `/data` in the pod)
- `DB_PATH` (default: `/data/nerdy_link_manager.db`)

After installation, the dashboard is available on:
- `http://localhost:$FRONTEND_PORT`

> [!IMPORTANT]
> Don't forget to change your password in the UI after the app is deployed.
