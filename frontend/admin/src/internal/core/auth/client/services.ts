'use client';

import { http, SingleResponse } from '@/internal/core/http';
import { JWT, JWTSchema } from '@/internal/core/auth/shared';

/**
 * Performs a login request to the server using the provided email and password.
 *
 * @param email - The email of the user attempting to log in.
 * @param password - The password of the user attempting to log in.
 * @returns A promise that resolves to a {@link SingleResponse} containing the JWT token upon successful login.
 */
export async function loginByCredentials(
  email: string,
  password: string,
): Promise<SingleResponse<JWT>> {
  const response = await http
    .post('api/v1/login', {
      json: {
        email,
        password,
      },
    })
    .json();

  return SingleResponse(JWTSchema).parseAsync(response);
}
