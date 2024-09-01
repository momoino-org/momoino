'use server';

import ky from 'ky';
import { isServer } from '../config/config';
import { cookies } from 'next/headers';

/**
 * Represents a generic response from the backend API.
 */
export type Response<T = undefined> = T extends undefined
  ? {
      message: string;
      messageId: string;
      timestamp: string;
    }
  : {
      message: string;
      messageId: string;
      timestamp: string;
      data: T;
    };

/**
 * Represents a paginated response from the backend API.
 */
export type PaginatedResponse<T> = Response<T> & {
  page: number;
  pageSize: number;
  totalRows: number;
  totalPages: number;
};

/**
 * Creates a configured instance of the Ky HTTP client with custom settings.
 * This instance is used for making requests to the backend API.
 *
 * @returns A configured instance of the Ky HTTP client.
 */
export const HTTPClient = ky.create({
  prefixUrl: isServer
    ? process.env.NEXT_BACKEND_HOST
    : process.env.NEXT_PUBLIC_BACKEND_HOST,
  throwHttpErrors: true,
  hooks: {
    beforeRequest: [
      async (request) => {
        if (isServer) {
          const accessToken = cookies().get('auth.token')?.value;
          request.headers.set('X-Auth-Access-Token', `Bearer ${accessToken}`);
        }
      },
    ],
  },
});
