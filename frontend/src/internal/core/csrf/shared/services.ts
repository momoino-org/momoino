import { http } from '@/internal/core/http';

export async function getNewCsrfToken(): Promise<{
  token: string;
  rawCookies: string[];
}> {
  const response = await http.get('api/v1/csrf-token');
  const csrfToken = response.headers.get('X-Csrf-Token');

  if (response.ok && csrfToken) {
    return {
      token: csrfToken,
      rawCookies: response.headers.getSetCookie(),
    };
  }

  throw new Error('Cannot get CSRF token', {
    cause: response,
  });
}
