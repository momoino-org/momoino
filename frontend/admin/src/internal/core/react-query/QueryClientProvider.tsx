'use client';

import { PropsWithChildren } from 'react';
import { QueryClientProvider as TanstackQueryClientProvider } from '@tanstack/react-query';
import { getQueryClient } from './utils';

/**
 * A React component that provides a Tanstack Query client to its children components.
 * It initializes a query client using the `getQueryClient` function and wraps the children
 * with the `TanstackQueryClientProvider` from `@tanstack/react-query`.
 *
 * @param props - The props for the component.
 * @param props.children - The children components to be wrapped with the query client provider.
 * @returns The component with the query client provider.
 */
export function QueryClientProvider(props: PropsWithChildren) {
  const queryClient = getQueryClient();

  return (
    <TanstackQueryClientProvider client={queryClient}>
      {props.children}
    </TanstackQueryClientProvider>
  );
}
