'use client';

import { AppContext } from '../core/ui';
import { isUseHttps } from '../core/config/config';
import {
  getUserProfile,
  isAccessTokenValid,
  loginByCredentials,
  setCookie,
} from './services';

export const AuthProvider: NonNullable<AppContext['authProvider']> = {
  login: (provider) => {
    switch (provider) {
      case 'credentials':
        return {
          mutationKey: ['login-crendetials'],
          mutationFn: (params) =>
            loginByCredentials(params.email, params.password),
          onSuccess: async (response) => {
            await setCookie({
              name: 'auth.token',
              value: response.data.accessToken,
              sameSite: 'strict',
              secure: isUseHttps,
              httpOnly: true,
            });

            localStorage.setItem(
              'auth.refreshToken',
              response.data.refreshToken,
            );
          },
        };
      default:
        throw new Error(`Unsupported login provider: ${provider}`);
    }
  },
  isAuthenticated: () => isAccessTokenValid(),
  profile: () => getUserProfile(),
};
