'use client';

import { Options } from 'ky';
import { http, SingleResponse } from '@/internal/core/http';
import { JWT, JWTSchema } from '@/internal/core/auth/shared';

/**
 * Performs a login request to the server using the provided username and password.
 *
 * @param username - The username of the user attempting to log in.
 * @param password - The password of the user attempting to log in.
 * @returns A promise that resolves to a {@link SingleResponse} containing the JWT token upon successful login.
 */
export async function loginByCredentials(
  username: string,
  password: string,
): Promise<SingleResponse<JWT>> {
  const response = await http
    .post('api/v1/login', {
      json: {
        username,
        password,
      },
    })
    .json();

  return SingleResponse(JWTSchema).parseAsync(response);
}

export async function refreshSession(options?: Options) {
  const response = await http.post('api/v1/token/renew', options).json();

  return SingleResponse(JWTSchema).parseAsync(response);
}
