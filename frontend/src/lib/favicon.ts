export function getFaviconCandidates(rawUrl: string): string[] {
  try {
    const u = new URL(rawUrl);
    const origin = u.origin;
    const host = u.hostname;

    return [
      `${origin}/favicon.ico`,
      `${origin}/apple-touch-icon.png`,
      `https://www.google.com/s2/favicons?domain=${host}&sz=64`,
      `https://icons.duckduckgo.com/ip3/${host}.ico`
    ];
  } catch {
    return [];
  }
}
