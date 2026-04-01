import { env } from '$env/dynamic/public';
import type { AuthResponse } from '$lib/./types';

const API_BASE = env.PUBLIC_API_URL;

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers = new Headers(options.headers || {});
  headers.set('Content-Type', 'application/json');

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }

  return (await res.json()) as T;
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
    })
};
