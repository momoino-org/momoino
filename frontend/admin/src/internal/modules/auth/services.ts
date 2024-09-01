'use server';

import { HTTPClient, Response } from '../core/http-client/http';
import { JWT, Profile } from '../core/ui/providers/AppProvider';
import { cookies } from 'next/headers';
import { ResponseCookie } from 'next/dist/compiled/@edge-runtime/cookies';
import { jwtDecode } from 'jwt-decode';
import { JwtPayload } from './type';
import { isEmpty } from 'radash';

/**
 * Performs a login request using the provided email and password.
 *
 * @param email The email of the user attempting to log in.
 * @param password The password of the user attempting to log in.
 *
 * @returns A Promise that resolves to a {@link Response} containing the JWT tokens.
 */
export async function loginByCredentials(
  email: string,
  password: string,
): Promise<Response<JWT>> {
  return HTTPClient.post('api/v1/login', {
    json: {
      email,
      password,
    },
  }).json<Response<JWT>>();
}

/**
 * Sets a cookie using the provided options.
 *
 * @param options The options for setting the cookie. This object should conform to the {@link ResponseCookie} interface.
 *
 * @returns A Promise that resolves to `void` once the cookie is set.
 */
export async function setCookie(options: ResponseCookie): Promise<void> {
  cookies().set(options);
}

/**
 * Checks if the access token is valid by making a request to the user profile endpoint.
 *
 * @returns A Promise that resolves to `true` if the access token is valid and the user profile
 * can be retrieved, or `false` if the access token is missing, invalid, or the user profile
 * endpoint returns an error.
 *
 * @remarks
 * This function sends a GET request to the 'api/v1/profile' endpoint using the {@link HTTPClient}.
 * If the request is successful (status code 200), the function resolves to `true`. If the request
 * fails (status code other than 200), the function resolves to `false`.
 *
 * The access token is extracted from the 'auth.token' cookie using the {@link cookies} function.
 * If the cookie is missing or its value is empty, the function resolves to `false`.
 */
export async function isAccessTokenValid(): Promise<boolean> {
  try {
    await HTTPClient.get('api/v1/profile').json<Response<Profile>>();
    return true;
  } catch (error) {
    return false;
  }
}

/**
 * Retrieves the user's profile information from the access token cookie.
 *
 * @returns A Promise that resolves to the user's profile information if the access token is valid,
 * or `null` if the access token is missing or invalid.
 *
 * @remarks
 * This function extracts the access token from the 'auth.token' cookie, and decodes it
 * to obtain the user's profile information. If the access token is missing or invalid, the function
 * returns `null`.
 */
export async function getUserProfile(): Promise<Profile | null> {
  const accessTokenCookie = cookies().get('auth.token');

  if (accessTokenCookie === undefined || isEmpty(accessTokenCookie.value)) {
    return null;
  }

  const payload = jwtDecode<JwtPayload>(accessTokenCookie.value);

  return {
    id: payload.sub,
    email: payload.email,
    firstName: payload.given_name,
    lastName: payload.family_name,
    username: payload.preferred_username,
    locale: payload.locale,
  };
}
