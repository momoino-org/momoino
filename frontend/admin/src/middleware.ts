'use server';

import { NextMiddleware, NextRequest, NextResponse } from 'next/server';
import { AuthProvider } from './internal/modules/auth/AuthProvider';
import { match } from 'path-to-regexp';
import { isEmpty } from 'radash';

const isSignInRoute = match('/auth/signin');
const privateRoutes = [match('/admin{/*path}')];

const signInMiddleware: NextMiddleware = async (request: NextRequest) => {
  const redirectTo = request.nextUrl.searchParams.get('redirectTo');

  if (isEmpty(redirectTo)) {
    request.nextUrl.searchParams.set('redirectTo', '/');
    return NextResponse.redirect(request.nextUrl);
  }

  // Don't need to call api get user profile if "auth.token" cookie does not exist
  const accessToken = request.cookies.get('auth.token');
  if (accessToken === undefined) {
    return NextResponse.next();
  }

  // If user is authenticated, redirect to the specified redirectTo URL
  if (await AuthProvider.isAuthenticated()) {
    const redirectURL = new URL(redirectTo!, request.url);
    return NextResponse.redirect(redirectURL);
  }

  // If user is not authenticated, redirect to the sign-in page
  return NextResponse.next();
};

const privateRouteMiddleware: NextMiddleware = async (request: NextRequest) => {
  const route = privateRoutes.find(
    (r) => r(request.nextUrl.pathname) !== false,
  );

  if (route) {
    const signInURL = new URL('/auth/signin', request.url);
    signInURL.searchParams.set('redirectTo', request.nextUrl.pathname);

    // Don't need to call api get user profile if "auth.token" cookie does not exist
    const accessToken = request.cookies.get('auth.token');
    if (accessToken === undefined) {
      return NextResponse.redirect(signInURL);
    }

    if (await AuthProvider.isAuthenticated()) {
      return NextResponse.next();
    }

    return NextResponse.redirect(signInURL);
  }

  return NextResponse.next();
};

export const middleware: NextMiddleware = async (request, event) => {
  if (isSignInRoute(request.nextUrl.pathname) !== false) {
    return signInMiddleware(request, event);
  }

  return privateRouteMiddleware(request, event);
};
