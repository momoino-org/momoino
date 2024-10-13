'use server';

import { cookies } from 'next/headers';
import { jwtDecode } from 'jwt-decode';
import { isEmpty } from 'radash';
import {
  JwtPayload,
  Profile,
  ProfileSchema,
} from '@/internal/core/auth/shared';
import { http } from '@/internal/core/http';

/**
 * Checks if the access token in the 'auth.token' cookie is valid by making a request to the server.
 *
 * @returns A Promise that resolves to `true` if the access token is valid, or `false` otherwise.
 */
export async function isAccessTokenValid(): Promise<boolean> {
  try {
    const accessTokenCookie = cookies().get('auth.token');

    if (accessTokenCookie === undefined || isEmpty(accessTokenCookie.value)) {
      return false;
    }

    await http
      .get('api/v1/profile', {
        headers: {
          'X-Auth-Access-Token': `Bearer ${accessTokenCookie.value}`,
        },
      })
      .json();
    return true;
  } catch (error) {
    return false;
  }
}

/**
 * Retrieves the user profile from the access token stored in the 'auth.token' cookie.
 *
 * @returns A Promise that resolves to the user profile if the access token is valid and present,
 * or `null` if the access token is missing or invalid.
 */
export async function getUserProfile(): Promise<Profile | null> {
  const accessTokenCookie = cookies().get('auth.token');

  if (accessTokenCookie === undefined || isEmpty(accessTokenCookie.value)) {
    return null;
  }

  const payload = jwtDecode<JwtPayload>(accessTokenCookie.value);

  return ProfileSchema.parseAsync({
    id: payload.sub,
    username: payload.preferred_username,
    email: payload.email,
    firstName: payload.given_name,
    lastName: payload.family_name,
    locale: payload.locale,
  });
}
