'use client';

import { useMutation, UseMutationOptions } from '@tanstack/react-query';
import { createContext, ProviderProps, useContext } from 'react';
import { Response } from '../../http-client/http';

export interface JWT {
  accessToken: string;
  refreshToken: string;
}

export interface Profile {
  id: string;
  username: string;
  email: string;
  firstName: string;
  lastName: string;
  locale: string;
}

/**
 * Defines the type for the AuthProvider interface, which provides methods for authentication-related operations.
 */
type AuthProvider = {
  /**
   * A method for logging in using a specified provider.
   *
   * @param provider - The name of the authentication provider (e.g., 'credentials', 'google', 'facebook').
   * @returns An object containing options for the useMutation hook, which handles the login mutation.
   *          The mutation returns a Response<JWT> on success and throws an Error on failure.
   *          The mutation accepts a Record<string, string> as input, representing the login credentials.
   */
  login(
    provider: string,
  ): UseMutationOptions<Response<JWT>, Error, Record<string, string>>;

  /**
   * A method for checking if the user is authenticated.
   *
   * @returns A Promise that resolves to a boolean value indicating whether the user is authenticated.
   */
  isAuthenticated(): Promise<boolean>;

  /**
   * A method for retrieving the user's profile information.
   *
   * @returns A Promise that resolves to a Profile object if the user is authenticated, or null if not.
   */
  profile(): Promise<Profile | null>;
};

/**
 * Defines the context type for the application, providing access to authentication-related methods.
 */
export type AppContext = {
  /**
   * An optional object containing authentication provider methods.
   */
  authProvider?: AuthProvider;
};

const AppContext = createContext<AppContext>({
  authProvider: undefined,
});

export const AppProvider = (props: ProviderProps<AppContext>) => {
  return <AppContext.Provider {...props} />;
};

/**
 * A custom React hook that handles the login process using the 'credentials' authentication provider.
 *
 * @returns A useMutation hook for handling the login mutation.
 *          The mutation returns a Response<JWT> on success and throws an Error on failure.
 *          The mutation accepts an object with 'email' and 'password' properties as input.
 *
 * @throws An Error if the authentication provider is not set up.
 */
export function useLoginByCredentials() {
  const appCtx = useContext(AppContext);

  if (appCtx.authProvider?.login) {
    return useMutation<
      Response<JWT>,
      Error,
      { email: string; password: string }
    >(appCtx.authProvider.login('credentials'));
  }

  throw new Error('Please setup the auth provider');
}
