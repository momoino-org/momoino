'use server';

import { cookies } from 'next/headers';
import { jwtDecode } from 'jwt-decode';
import { isEmpty } from 'radash';
import { U } from 'ts-toolbelt';
import { RequestCookie } from 'next/dist/compiled/@edge-runtime/cookies';
import { NextRequest, NextResponse } from 'next/server';
import { getCsrfFrom } from '@/internal/core/csrf/server';
import { http } from '@/internal/core/http';
import { setCookie } from '@/internal/core/utils/server';
import {
  CsrfHeaderName,
  IdentityCookie,
  JwtPayload,
  Profile,
  ProfileSchema,
  SessionCookie,
} from '@/internal/core/auth/shared';

/**
 * Retrieves the user profile from the access token stored in the 'auth.token' cookie.
 *
 * @returns A Promise that resolves to the user profile if the access token is valid and present,
 * or `null` if the access token is missing or invalid.
 */
export async function getUserProfile(): Promise<U.Nullable<Profile>> {
  const accessTokenCookie = cookies().get(IdentityCookie);

  if (accessTokenCookie === undefined || isEmpty(accessTokenCookie.value)) {
    return null;
  }

  const payload = jwtDecode<JwtPayload>(accessTokenCookie.value);

  return ProfileSchema.parseAsync({
    id: payload.sub,
    sid: payload.sid,
    username: payload.preferred_username,
    email: payload.email,
    firstName: payload.given_name,
    lastName: payload.family_name,
    locale: payload.locale,
  });
}

/**
 * Retrieves the identity cookie from a request. Throws an error if the cookie is missing or empty.
 *
 * @param request - The incoming request object containing cookies.
 * @returns The identity cookie as a RequestCookie object.
 *
 * @throws Error if the identity cookie is missing or has an empty value.
 */
export async function getIdentityCookie(
  request: NextRequest,
): Promise<RequestCookie> {
  const identityCookie = request.cookies.get(IdentityCookie);

  if (!identityCookie || isEmpty(identityCookie.value)) {
    throw new Error('Missing identity cookie');
  }

  return identityCookie;
}

/**
 * Retrieves the session cookie from a request. Throws an error if the cookie is missing or empty.
 *
 * @param request - The incoming request object containing cookies.
 * @returns The session cookie as a RequestCookie object.
 *
 * @throws Error if the session cookie is missing or has an empty value.
 */
export async function getSessionCookie(
  request: NextRequest,
): Promise<RequestCookie> {
  const sessionCookie = request.cookies.get(SessionCookie);

  if (!sessionCookie || isEmpty(sessionCookie.value)) {
    throw new Error('Missing session cookie');
  }

  return sessionCookie;
}

/**
 * Injects a new access token into the response by making a POST request to the renew token endpoint.
 * The function retrieves the session and CSRF cookies from the request, and sends a POST request
 * to the renew token endpoint with the necessary headers. It then sets the new access token and CSRF
 * cookie in the response.
 *
 * @param request - The incoming NextRequest object containing the session and CSRF cookies.
 * @param response - The NextResponse object to which the new access token and CSRF cookie will be added.
 * @returns A Promise that resolves to the updated NextResponse object with the new access token and CSRF cookie.
 *
 * @throws Error if the session or CSRF cookies are missing or have empty values.
 */
export async function injectNewAccessToken(
  request: NextRequest,
  response: NextResponse,
): Promise<NextResponse> {
  const sessionCookie = await getSessionCookie(request);
  const csrf = await getCsrfFrom(response);

  const refreshSessionResponse = await http.post('api/v1/token/renew', {
    headers: {
      Cookie: `${csrf.cookie.name}=${csrf.cookie.value}; ${sessionCookie.name}=${sessionCookie.value}`,
      [CsrfHeaderName]: csrf.token,
    },
  });

  for (const rawCookie of refreshSessionResponse.headers.getSetCookie()) {
    setCookie(response, rawCookie);
  }

  return response;
}

export async function isAccessTokenValid(
  request: NextRequest,
): Promise<boolean> {
  try {
    const identityCookie = await getIdentityCookie(request);
    const sessionCookie = await getSessionCookie(request);
    const payload = jwtDecode<JwtPayload>(identityCookie.value);

    return payload.sid === sessionCookie.value;
  } catch (err) {
    return false;
  }
}

export async function canRenewAccessToken(
  request: NextRequest,
): Promise<boolean> {
  try {
    await getSessionCookie(request);
    return true;
  } catch (err) {
    return false;
  }
}
