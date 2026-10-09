import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AuthenticationRequiredError, createApiClient } from './client';

describe('API client', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('sends only current Firebase ID token as bearer without cookies', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(null, { status: 204 }),
    );
    const client = createApiClient(
      async () => 'firebase-id-token',
      'https://api.example.test',
    );

    await client.request('/api/v1/me');

    expect(fetchMock).toHaveBeenCalledWith(
      'https://api.example.test/api/v1/me',
      expect.objectContaining({
        credentials: 'omit',
        headers: expect.any(Headers),
      }),
    );
    const request = fetchMock.mock.calls[0]?.[1];
    expect(new Headers(request?.headers).get('Authorization')).toBe(
      'Bearer firebase-id-token',
    );
  });

  it('rejects requests without an ID token or HTTPS outside local', async () => {
    const unauthenticatedClient = createApiClient(async () => null, 'https://api.example.test');

    await expect(unauthenticatedClient.request('/api/v1/me'))
      .rejects.toBeInstanceOf(AuthenticationRequiredError);

    const insecureClient = createApiClient(async () => 'token', 'http://api.example.test');
    await expect(insecureClient.request('/api/v1/me'))
      .rejects.toThrow('API requests require HTTPS outside Local.');
  });

  it('requires explicit API URL in production builds', async () => {
    const client = createApiClient(async () => 'token', '');

    await expect(client.request('/api/v1/me'))
      .rejects.toThrow('API base URL is not configured for this build.');
  });

  it('rejects protocol-relative paths that could redirect bearer traffic', async () => {
    const client = createApiClient(async () => 'token', 'https://api.example.test');

    await expect(client.request('//attacker.example.test/collect'))
      .rejects.toThrow('API client accepts same-origin relative paths only.');
  });

  it('rejects backslash paths that browsers normalize to another origin', async () => {
    const client = createApiClient(async () => 'token', 'https://api.example.test');

    await expect(client.request('/\\attacker.example.test/collect'))
      .rejects.toThrow('API client accepts same-origin relative paths only.');
  });
});
