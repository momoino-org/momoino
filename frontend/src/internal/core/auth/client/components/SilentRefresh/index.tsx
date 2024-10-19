'use client';

import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { refreshSession } from '../../services';
import { authSlice } from '../../store';
import { getUserProfile } from '../../../server';
import { useAppDispatch } from '@/internal/core/ui';
import { notification } from '@/internal/core/ui/toast';

export function SilentRefresh() {
  const dispatch = useAppDispatch();
  const queryClient = useQueryClient();
  const router = useRouter();

  const { data, isError, error } = useQuery({
    enabled: true,
    queryKey: ['renew-access-token'],
    refetchInterval: 50_000,
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
    refetchOnReconnect: true,
    retry: 2,
    queryFn: ({ signal }) => refreshSession({ signal }),
  });

  useEffect(() => {
    return () => {
      notification.toast({
        severity: 'error',
        message: 'Your session is expired.',
      });
    };
  }, []);

  useEffect(() => {
    if (isError) {
      router.refresh();
    } else {
      getUserProfile().then((userProfile) => {
        dispatch(authSlice.actions.setProfile(userProfile));
      });
    }

    return () => {
      queryClient.cancelQueries({
        queryKey: ['renew-access-token'],
      });
    };
  }, [data, error, dispatch, queryClient, router, isError]);

  return null;
}
