'use client';

import { combineSlices, configureStore } from '@reduxjs/toolkit';
import { Provider, useDispatch, useSelector, useStore } from 'react-redux';
import { PropsWithChildren, useRef } from 'react';
import { isProduction } from '@/internal/core/config';

export interface LazyLoadedSlices {}

export const rootReducer =
  combineSlices().withLazyLoadedSlices<LazyLoadedSlices>();

const makeStore = () => {
  const store = configureStore({
    reducer: rootReducer,
    devTools: !isProduction,
    middleware: (getDefaultMiddleware) => getDefaultMiddleware(),
  });

  return store;
};

export type AppStore = ReturnType<typeof makeStore>;
export type RootState = ReturnType<AppStore['getState']>;
export type AppDispatch = AppStore['dispatch'];

export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
export const useAppSelector = useSelector.withTypes<RootState>();
export const useAppStore = useStore.withTypes<AppStore>();

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
export function StoreProvider(props: PropsWithChildren) {
  const storeRef = useRef<AppStore>(makeStore());

  return <Provider store={storeRef.current}>{props.children}</Provider>;
}
