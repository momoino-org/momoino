'use server';

import {
  RequestCookie,
  ResponseCookie,
} from 'next/dist/compiled/@edge-runtime/cookies';
import { NextRequest, NextResponse } from 'next/server';
import { isEmpty } from 'radash';
import { getNewCsrfToken } from '@/internal/core/csrf/shared';
import { CsrfCookie, CsrfHeaderName } from '@/internal/core/auth/shared';
import { setCookie } from '@/internal/core/utils/server';

/**
 * Retrieves the CSRF cookie and CSRF token header from the given request or response.
 * Throws an error if either the CSRF cookie or token is missing or empty.
 *
 * @param source - The request or response object containing the cookies and headers.
 * @returns An object containing the CSRF cookie and token.
 *
 * @throws Error if the CSRF cookie or CSRF token is missing or has an empty value.
 */
export async function getCsrfFrom(source: NextRequest | NextResponse): Promise<{
  cookie: RequestCookie | ResponseCookie;
  token: string;
}> {
  const csrfCookie = source.cookies.get(CsrfCookie);
  const csrfToken = source.headers.get(CsrfHeaderName);

  if (!csrfCookie || isEmpty(csrfCookie.value)) {
    throw new Error('Missing CSRF cookie');
  }

  if (!csrfToken || isEmpty(csrfToken)) {
    throw new Error('Missing CSRF token');
  }

  return {
    cookie: csrfCookie,
    token: csrfToken,
  };
}

export async function injectCsrfToken(
  response: NextResponse,
): Promise<NextResponse> {
  const { rawCookies, token } = await getNewCsrfToken();

  response.headers.set(CsrfHeaderName, token);

  for (const rawCookie of rawCookies) {
    setCookie(response, rawCookie);
  }

  return response;
}
