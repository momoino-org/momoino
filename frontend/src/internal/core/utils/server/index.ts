'use server';

import { ResponseCookie } from 'next/dist/compiled/@edge-runtime/cookies';
import { NextResponse } from 'next/server';
import { parseString } from 'set-cookie-parser';

/**
 * Sets a cookie in the response based on the provided raw cookie string.
 *
 * @param response - The NextResponse object to which the cookie will be added.
 * @param rawCookie - The raw cookie string to be parsed and added to the response.
 */
export async function setCookie(
  response: NextResponse,
  rawCookie: string,
): Promise<void> {
  const parsedCookie = parseString(rawCookie);

  response.cookies.set({
    name: parsedCookie.name,
    value: parsedCookie.value,
    path: parsedCookie.path,
    expires: parsedCookie.expires,
    maxAge: parsedCookie.maxAge,
    domain: parsedCookie.domain,
    secure: parsedCookie.secure,
    httpOnly: parsedCookie.httpOnly,
    sameSite: parsedCookie.sameSite as ResponseCookie['sameSite'],
  });
}
