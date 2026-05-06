const CACHE_KEY = 'favicon-cache-v1';

type CacheMap = Record<string, string>;

function loadCache(): CacheMap {
  if (typeof localStorage === 'undefined') return {};
  try {
    return JSON.parse(localStorage.getItem(CACHE_KEY) || '{}');
  } catch {
    return {};
  }
}

function saveCache(cache: CacheMap) {
  if (typeof localStorage === 'undefined') return;
  localStorage.setItem(CACHE_KEY, JSON.stringify(cache));
}

function getHost(rawUrl: string): string | null {
  try {
    return new URL(rawUrl).hostname;
  } catch {
    return null;
  }
}

export function getCachedFavicon(rawUrl: string): string | null {
  const host = getHost(rawUrl);
  if (!host) return null;
  const cache = loadCache();
  return cache[host] ?? null;
}

export function setCachedFavicon(rawUrl: string, faviconUrl: string) {
  const host = getHost(rawUrl);
  if (!host) return;
  const cache = loadCache();
  cache[host] = faviconUrl;
  saveCache(cache);
}

export function getFaviconCandidates(rawUrl: string): string[] {
  try {
    const u = new URL(rawUrl);
    const origin = u.origin;
    const host = u.hostname;

    return [
      `${origin}/favicon.ico`,
      `${origin}/favicon.png`,
      `${origin}/favicon.svg`,
      `${origin}/apple-touch-icon.png`,
      `${origin}/apple-touch-icon-precomposed.png`,
      `${origin}/android-chrome-192x192.png`,
      `https://www.google.com/s2/favicons?domain_url=${origin}&sz=64`,
      `https://icons.duckduckgo.com/ip3/${host}.ico`
    ];
  } catch {
    return [];
  }
}
