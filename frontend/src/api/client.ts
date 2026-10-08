const defaultApiBaseUrl = import.meta.env.VITE_API_BASE_URL?.trim() ||
  (import.meta.env.DEV ? 'http://localhost:5080' : '');

export class AuthenticationRequiredError extends Error {
  public constructor() {
    super('A signed-in Firebase user is required.');
    this.name = 'AuthenticationRequiredError';
  }
}

export type IdTokenProvider = (forceRefresh?: boolean) => Promise<string | null>;

function isLocalHost(hostname: string): boolean {
  return hostname === 'localhost' || hostname === '127.0.0.1' || hostname === '[::1]' || hostname === '::1';
}

function resolveApiUrl(path: string, baseUrl: string): string {
  if (!path.startsWith('/') || path.startsWith('//')) {
    throw new Error('API client accepts same-origin relative paths only.');
  }

  if (!baseUrl) {
    throw new Error('API base URL is not configured for this build.');
  }

  const base = new URL(baseUrl);
  const url = new URL(path, base);
  if (base.username || base.password || url.origin !== base.origin) {
    throw new Error('API client accepts same-origin relative paths only.');
  }

  if (url.protocol !== 'https:' && !(import.meta.env.DEV && isLocalHost(url.hostname))) {
    throw new Error('API requests require HTTPS outside Local.');
  }

  return url.toString();
}

export function createApiClient(getIdToken: IdTokenProvider, baseUrl = defaultApiBaseUrl) {
  return {
    async request(path: string, init: RequestInit = {}): Promise<Response> {
      const idToken = await getIdToken();
      if (!idToken) {
        throw new AuthenticationRequiredError();
      }

      const headers = new Headers(init.headers);
      headers.set('Authorization', `Bearer ${idToken}`);

      return fetch(resolveApiUrl(path, baseUrl), {
        ...init,
        credentials: 'omit',
        headers,
      });
    },
  };
}
