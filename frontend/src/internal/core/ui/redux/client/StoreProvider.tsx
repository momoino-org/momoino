'use client';

import { PropsWithChildren, useRef } from 'react';
import { Provider } from 'react-redux';
import { U } from 'ts-toolbelt';
import { AppStore, makeStore } from '..';
import { authSlice } from '@/internal/core/auth/client/store';
import { Profile } from '@/internal/core/auth/shared';

interface StoreProps {
  profile: U.Nullable<Profile>;
}

/**
 * StoreProvider component that initializes a global store and provides it to
 * the React component tree using the Context API.
 *
 * This component is responsible for creating a single instance of the application
 * store using the `makeStore()` function, which will be shared across all child
 * components via the React Context `Provider`. It ensures that the store is created
 * only once, even if the component is re-rendered multiple times.
 *
 * @param props - The props object which contains the child components
 * to be wrapped by the Provider. This enables passing down the global store to all
 * components nested within the StoreProvider.
 *
 * @returns The Provider component that supplies the store to the rest of the application
 * via React's Context API.
 */
export function StoreProvider(props: PropsWithChildren<StoreProps>) {
  const storeRef = useRef<AppStore | null>(null);

  if (!storeRef.current) {
    storeRef.current = makeStore();
    storeRef.current.dispatch(authSlice.actions.setProfile(props.profile));
  }

  return <Provider store={storeRef.current}>{props.children}</Provider>;
}
