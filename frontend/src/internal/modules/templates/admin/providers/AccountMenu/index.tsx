'use client';

import { useSession } from 'next-auth/react';
import { AccountMenu } from '../../components/AccountMenu';

export function AXccountMenu() {
  const session = useSession();

  return (
    <>
      <AccountMenu
        slotProps={{
          avatar: {
            src: `https://ui-avatars.com/api/?rounded=true&name=${session.data?.user?.firstName}&size=24`,
          },
          tooltip: {
            title: 'Account',
          },
        }}
      />
    </>
  );
}
