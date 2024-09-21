import ky from 'ky';
import { isServer } from '@/internal/core/config';

/**
 * Creates a configured instance of the Ky HTTP client with custom settings.
 * This instance is used for making requests to the backend API.
 *
 * @returns A configured instance of the Ky HTTP client.
 */
export const http = ky.create({
  prefixUrl: isServer
    ? process.env.NEXT_BACKEND_HOST
    : process.env.NEXT_PUBLIC_BACKEND_HOST,
  throwHttpErrors: true,
  credentials: 'include',
});
