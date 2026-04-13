import { env } from '$env/dynamic/public';
import type {
  AuthResponse,
  BackendLinkItem,
  LinkItem,
  CreateLinkPayload,
  UpdateLinkPayload
} from '$lib/types';

const API_BASE = env.PUBLIC_API_URL;

function getToken() {
  return typeof window !== 'undefined' ? localStorage.getItem('token') : null;
}

function mapLinkFromBackend(item: BackendLinkItem): LinkItem {
  return {
    id: String(item.id),
    name: item.title,
    url: item.url,
    description: item.description ?? null,
    createdAt: item.created_at,
    updatedAt: item.updated_at,
    clicks: item.click_count
  };
}

function toBackendPayload(payload: CreateLinkPayload | UpdateLinkPayload) {
  return {
    title: payload.name,
    url: payload.url,
    description: payload.description ?? null
  };
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers = new Headers(options.headers || {});
  headers.set('Content-Type', 'application/json');

  const token = getToken();
  if (token) headers.set('Authorization', `Bearer ${token}`);

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }

  const contentType = res.headers.get('content-type') ?? '';
  if (contentType.includes('application/json')) return (await res.json()) as T;
  return '' as T;
}

export const api = {
  register: (email: string, password: string) =>
    request<AuthResponse>('/api/register', {
      method: 'POST',
      body: JSON.stringify({ email, password })
    }),

  login: (email: string, password: string) =>
    request<AuthResponse>('/api/login', {
      method: 'POST',
      body: JSON.stringify({ email, password })
    }),

  getLinks: async (withStats = true): Promise<LinkItem[]> => {
    const raw = await request<BackendLinkItem[]>(`/api/links${withStats ? '?stats=true' : ''}`);
    return raw.map(mapLinkFromBackend);
  },

  registerClick: (id: string) =>
    request<void>(`/api/links/${id}/click`, { method: 'POST' }),

  createLink: async (payload: CreateLinkPayload): Promise<LinkItem> => {
    const raw = await request<BackendLinkItem>('/api/links?stats=true', {
      method: 'POST',
      body: JSON.stringify(toBackendPayload(payload))
    });
    return mapLinkFromBackend(raw);
  },

  updateLink: async (id: string, payload: UpdateLinkPayload): Promise<LinkItem> => {
    const raw = await request<BackendLinkItem>(`/api/links/${id}?stats=true`, {
      method: 'PUT',
      body: JSON.stringify(toBackendPayload(payload))
    });
    return mapLinkFromBackend(raw);
  },

  deleteLink: (id: string) =>
    request<void>(`/api/links/${id}`, {
      method: 'DELETE'
    })
};
