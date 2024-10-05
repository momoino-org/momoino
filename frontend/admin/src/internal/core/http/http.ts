import ky from 'ky';
import { backendOrigin } from '@/internal/core/config';

/**
 * Creates a configured instance of the Ky HTTP client with custom settings.
 * This instance is used for making requests to the backend API.
 *
 * @returns A configured instance of the Ky HTTP client.
 */
export const http = ky.create({
  prefixUrl: backendOrigin,
  throwHttpErrors: true,
  credentials: 'include',
  retry: 0,
});
