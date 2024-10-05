import { ReadonlyURLSearchParams } from 'next/navigation';
import { useCallback } from 'react';
import { AuthenticationMessageData } from '../useOAuth2Listener';

export function useOAuth2Callback() {
  const validateOAuthState = (state: string | null) => {
    if (state === null) {
      throw new Error(
        'No OAuth state provided. Please try again, and if the issue persists, contact the system administrator for assistance.',
      );
    }

    if (sessionStorage.getItem('auth.state') !== state) {
      throw new Error(
        'Invalid OAuth state. Please try again, and if the issue persists, contact the system administrator for assistance.',
      );
    }
  };

  const getCodeVerifier = () => {
    const codeVerifier = sessionStorage.getItem('auth.verifier');

    if (!codeVerifier) {
      throw new Error(
        'Unable to retrieve the code verifier. Please try again, and if the issue persists, contact the system administrator for assistance.',
      );
    }

    return codeVerifier;
  };

  const isUsePkce = () => sessionStorage.getItem('auth.usePkce') === 'true';

  const getProvider = () => sessionStorage.getItem('auth.provider');

  const getBackendCallbackURL = useCallback(
    (searchParams: ReadonlyURLSearchParams) => {
      const params = new URLSearchParams(searchParams);

      if (isUsePkce()) {
        params.set('verifier', getCodeVerifier());
      }

      return `api/v1/login/providers/${getProvider()}/callback?${params.toString()}`;
    },
    [],
  );

  const success = useCallback(() => {
    (window.opener as WindowProxy).postMessage({
      source: `useOAuth2`,
      payload: {
        status: 'success',
      },
    } satisfies AuthenticationMessageData);
  }, []);

  const fail = useCallback((error: unknown) => {
    (window.opener as WindowProxy).postMessage({
      source: `useOAuth2`,
      payload: {
        status: 'error',
        details: error,
      },
    } satisfies AuthenticationMessageData);
  }, []);

  return {
    validateOAuthState,
    getBackendCallbackURL,
    success,
    fail,
  };
}
