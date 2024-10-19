'use server';

import { NextMiddleware, NextResponse } from 'next/server';
import { chain } from 'radash';
import { frontendOrigin } from './internal/core/config';
import { injectCsrfToken } from '@/internal/core/csrf/server';
import {
  canRenewAccessToken,
  injectNewAccessToken,
  isAccessTokenValid,
} from '@/internal/core/auth/server';

const homeURL = new URL('/', frontendOrigin);
const signInURL = new URL('/auth/signin', frontendOrigin);

const signInMiddleware: NextMiddleware = async (request) => {
  if (await isAccessTokenValid(request)) {
    return chain(NextResponse.redirect, injectCsrfToken)(homeURL);
  }

  try {
    if (await canRenewAccessToken(request)) {
      return chain(NextResponse.redirect, injectCsrfToken, async (response) =>
        injectNewAccessToken(request, await response),
      )(homeURL);
    }
  } catch (err) {
    return chain(NextResponse.next, injectCsrfToken)();
  }

  return chain(NextResponse.next, injectCsrfToken)();
};

const protectedRouteMiddleware: NextMiddleware = async (request) => {
  if (await isAccessTokenValid(request)) {
    return chain(NextResponse.next, injectCsrfToken)();
  }

  try {
    if (await canRenewAccessToken(request)) {
      return chain(NextResponse.next, injectCsrfToken, async (response) =>
        injectNewAccessToken(request, await response),
      )();
    }
  } catch (err) {
    return chain(NextResponse.redirect, injectCsrfToken)(signInURL);
  }

  return chain(NextResponse.redirect, injectCsrfToken)(signInURL);
};

export const middleware: NextMiddleware = async (request, event) => {
  const { pathname } = request.nextUrl;

  if (pathname === '/auth/signin') {
    return signInMiddleware(request, event);
  }

  return protectedRouteMiddleware(request, event);
};

export const config = {
  matcher: [
    {
      source: '/auth/signin',
      missing: [
        { type: 'header', key: 'next-router-prefetch' },
        { type: 'header', key: 'next-action' },
        { type: 'header', key: 'purpose', value: 'prefetch' },
      ],
    },
    {
      source: '/admin/:path*',
      missing: [
        { type: 'header', key: 'next-router-prefetch' },
        { type: 'header', key: 'next-action' },
        { type: 'header', key: 'purpose', value: 'prefetch' },
      ],
    },
  ],
};
