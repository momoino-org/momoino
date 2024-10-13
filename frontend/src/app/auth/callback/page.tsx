'use client';

import { useQuery } from '@tanstack/react-query';
import { useSearchParams } from 'next/navigation';
import { useEffect } from 'react';
import { http } from '@/internal/core/http';
import { useOAuth2Callback } from '@/internal/core/auth/client';

export default function OAuth2CallbackPage() {
  const searchParams = useSearchParams();

  const { validateOAuthState, success, fail, getBackendCallbackURL } =
    useOAuth2Callback();

  const { isSuccess, isFetched, isError, error } = useQuery({
    queryKey: ['oauth2'],
    queryFn: () => {
      validateOAuthState(searchParams.get('state'));
      return http.get(getBackendCallbackURL(searchParams));
    },
  });

  useEffect(() => {
    if (isFetched) {
      if (isError) {
        fail(error);
      } else {
        success();
      }

      window.close();
    }
  }, [isSuccess, isError, fail, error, success, isFetched]);

  return 'Loading...';
}
