import NextAuth, { User } from 'next-auth';
import Keycloak from 'next-auth/providers/keycloak';
import { JWT } from 'next-auth/jwt';
import { decodeJwt } from 'jose';
import { ProfileSchema } from './internal/core/auth/shared';

declare module 'next-auth' {
  interface User {
    username?: string;
    emailVerified?: boolean;
    firstName?: string;
    lastName?: string;
    locale?: string;
  }
}

declare module 'next-auth/jwt' {
  interface JWT {
    accessToken: string;
    accessTokenExpiresAt: number;
    refreshToken: string;
    refreshTokenExpiresAt: number;
    user: User;
  }
}

function isExpired(deadline: number): boolean {
  const now = Date.now();
  console.log({ now, deadline: deadline, isExpired: now > deadline });
  return now > deadline;
}

export const { auth, handlers, signIn, signOut } = NextAuth({
  providers: [Keycloak],

  callbacks: {
    async jwt({ token, profile, account }) {
      if (account) {
        return {
          user: await ProfileSchema.parseAsync({
            id: profile?.sub,
            username: profile?.preferred_username,
            email: profile?.email,
            firstName: profile?.given_name,
            lastName: profile?.family_name,
            locale: profile?.locale,
            emailVerified: profile?.email_verified,
          }),
          accessToken: account.access_token!,
          accessTokenExpiresAt: account.expires_at! * 1000,
          refreshToken: account.refresh_token!,
          refreshTokenExpiresAt:
            Date.now() + Number(account.refresh_expires_in!) * 1000,
        } satisfies JWT;
      } else if (!isExpired(token.accessTokenExpiresAt)) {
        return token;
      } else {
        if (isExpired(token.refreshTokenExpiresAt as number)) {
          throw new Error('Refresh token is expired');
        }

        console.debug('[START] Renew-ing access token using refresh token');
        console.debug('[OLD] Refresh token: ' + token.refreshToken);
        const request = await fetch(
          `${process.env.AUTH_KEYCLOAK_ISSUER}/protocol/openid-connect/token`,
          {
            method: 'post',
            headers: {
              'Content-Type': 'application/x-www-form-urlencoded',
            },
            cache: 'no-cache',
            body: new URLSearchParams({
              client_id: String(process.env.AUTH_KEYCLOAK_ID),
              client_secret: String(process.env.AUTH_KEYCLOAK_SECRET),
              grant_type: 'refresh_token',
              refresh_token: String(token.refreshToken),
            }),
          },
        );

        const response = await request.json();
        if (!request.ok) {
          console.debug(
            '[END] Renew-ing access token using refresh token failed',
            response,
          );
          throw response;
        }
        console.debug('[NEW] Refresh token: ' + response.refresh_token);

        const t = decodeJwt(response.access_token);
        return {
          user: await ProfileSchema.parseAsync({
            id: t.sub,
            username: t?.preferred_username,
            email: t?.email,
            firstName: t?.given_name,
            lastName: t?.family_name,
            locale: t?.locale,
            emailVerified: t?.email_verified,
          }),
          accessToken: response.access_token!,
          accessTokenExpiresAt:
            Date.now() + Number(response.expires_in!) * 1000,
          refreshToken: response.refresh_token!,
          refreshTokenExpiresAt:
            Date.now() + Number(response.refresh_expires_in!) * 1000,
        } satisfies JWT;
      }
    },

    async session({ session, token }) {
      return {
        ...session,
        user: token.user,
      };
    },
  },
});
