import ky from 'ky';
import { backendOrigin, isClient } from '@/internal/core/config';

/**
 * A list of HTTP methods considered "safe" for requests.
 *
 * Safe methods are defined as those that do not alter the state of the server
 * and are intended for data retrieval or non-modifying operations.
 * These methods can be called without causing side effects and are generally
 * safe to be cached and retried.
 *
 * The current safe HTTP methods are:
 * - `GET`: Retrieve data from the server.
 * - `HEAD`: Similar to GET but only retrieves the headers.
 * - `OPTIONS`: Describe the communication options for the target resource.
 */
const safeHttpMethods: ReadonlyArray<string> = ['GET', 'HEAD', 'OPTIONS'];

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
  headers: {
    'X-Requested-With': process.env.NEXT_PUBLIC_APP_NAME,
  },
  hooks: {
    beforeRequest: [
      (request) => {
        if (isClient && !safeHttpMethods.includes(request.method)) {
          const csrfToken = document.querySelector<HTMLMetaElement>(
            'meta[name="csrf-token"]',
          );

          if (csrfToken) {
            request.headers.set('X-Csrf-Token', csrfToken.content);
          }
        }
      },
    ],
  },
});
