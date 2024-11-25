import { combineSlices, configureStore } from '@reduxjs/toolkit';
import { isProduction } from '@/internal/core/config';

export interface LazyLoadedSlices {}

export const rootReducer = combineSlices(
  {},
).withLazyLoadedSlices<LazyLoadedSlices>();

export const makeStore = () => {
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
