import { env } from '$env/dynamic/public';
import type { AuthResponse, BackendLinkItem, LinkItem } from '$lib/types';

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
    request<AuthResponse>('/register', {
      method: 'POST',
      body: JSON.stringify({ email, password })
    }),

  login: (email: string, password: string) =>
    request<AuthResponse>('/login', {
      method: 'POST',
      body: JSON.stringify({ email, password })
    }),

  getLinks: async (withStats = true): Promise<LinkItem[]> => {
    const raw = await request<BackendLinkItem[]>(`/links${withStats ? '?stats=true' : ''}`);
    return raw.map(mapLinkFromBackend);
  },

  registerClick: (id: string) =>
    request<void>(`/links/${id}/click`, {
      method: 'POST'
    })
};
