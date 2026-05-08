const FAVICON_CACHE_KEY = 'favicon-cache-v1';
const CUSTOM_ICON_CACHE_KEY = 'custom-link-icons-v1';

type CacheMap = Record<string, string>;

function loadCache(key: string): CacheMap {
  if (typeof localStorage === 'undefined') return {};
  try {
    return JSON.parse(localStorage.getItem(key) || '{}');
  } catch {
    return {};
  }
}

function saveCache(key: string, cache: CacheMap) {
  if (typeof localStorage === 'undefined') return;
  localStorage.setItem(key, JSON.stringify(cache));
}

export function getCachedFavicon(rawUrl: string): string | null {
  try {
    const host = new URL(rawUrl).hostname;
    const cache = loadCache(FAVICON_CACHE_KEY);
    return cache[host] ?? null;
  } catch {
    return null;
  }
}

export function setCachedFavicon(rawUrl: string, faviconUrl: string) {
  try {
    const host = new URL(rawUrl).hostname;
    const cache = loadCache(FAVICON_CACHE_KEY);
    cache[host] = faviconUrl;
    saveCache(FAVICON_CACHE_KEY, cache);
  } catch { }
}

export function getCustomIcon(linkId: string): string | null {
  const cache = loadCache(CUSTOM_ICON_CACHE_KEY);
  return cache[linkId] ?? null;
}

export function setCustomIcon(linkId: string, dataUrl: string) {
  const cache = loadCache(CUSTOM_ICON_CACHE_KEY);
  cache[linkId] = dataUrl;
  saveCache(CUSTOM_ICON_CACHE_KEY, cache);
}

export function removeCustomIcon(linkId: string) {
  const cache = loadCache(CUSTOM_ICON_CACHE_KEY);
  delete cache[linkId];
  saveCache(CUSTOM_ICON_CACHE_KEY, cache);
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
