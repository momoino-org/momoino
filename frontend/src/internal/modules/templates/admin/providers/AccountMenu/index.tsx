'use client';

import { AccountMenu } from '../../components/AccountMenu';
import { SilentRefresh } from '@/internal/core/auth/client/components/SilentRefresh';
import { useAppSelector } from '@/internal/core/ui';

export function AXccountMenu() {
  const userProfile = useAppSelector((state) => state.auth.profile);

  if (!userProfile) {
    return null;
  }

  return (
    <>
      <SilentRefresh />
      <AccountMenu
        slotProps={{
          avatar: {
            src: `https://ui-avatars.com/api/?rounded=true&name=${userProfile.firstName}&size=24`,
          },
          tooltip: {
            title: 'Account',
          },
        }}
      />
    </>
  );
}
