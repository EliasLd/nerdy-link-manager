import { env } from '$env/dynamic/public';
import type { AuthResponse, LinkItem } from '$lib/./types';

const API_BASE = env.PUBLIC_API_URL;

function getToken() {
  return typeof window !== 'undefined' ? localStorage.getItem('token') : null;
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

  getLinks: (withStats = true) =>
    request<LinkItem[]>(`/links${withStats ? '?stats=true' : ''}`),

  registerClick: (id: string) =>
    request<void>(`/links/${id}/click`, {
      method: 'POST'
    })
};
