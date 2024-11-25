import {
  QueryClient,
  defaultShouldDehydrateQuery,
  isServer,
} from '@tanstack/react-query';
import { HTTPError } from 'ky';
import { notification } from '@/internal/core/ui';

/**
 * A global variable to hold the browser-side QueryClient instance.
 */
let browserQueryClient: QueryClient | undefined = undefined;

/**
 * Creates a new QueryClient instance with default options.
 *
 * @returns A new QueryClient instance.
 */
function makeQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60 * 1000, // Cache data for 1 minute
        refetchOnWindowFocus: false,
        retry: false,
        throwOnError(error) {
          if (error instanceof HTTPError) {
            return false;
          }

          return true;
        },
      },
      mutations: {
        async onError(error) {
          if (error instanceof HTTPError) {
            const response = await error.response.json();
            notification.toast({
              severity: 'error',
              message: response.message,
            });
          }
        },
        throwOnError(error) {
          if (error instanceof HTTPError) {
            return false;
          }

          return true;
        },
      },
      dehydrate: {
        // Include pending queries in dehydration
        shouldDehydrateQuery: (query) =>
          defaultShouldDehydrateQuery(query) ||
          query.state.status === 'pending',
      },
    },
  });
}

/**
 * Returns a QueryClient instance based on the environment.
 *
 * @returns A QueryClient instance for the server or browser.
 */
export function getQueryClient() {
  if (isServer) {
    // Server: always make a new query client
    return makeQueryClient();
  } else {
    // Browser: make a new query client if we don't already have one
    // This is very important, so we don't re-make a new client if React
    // suspends during the initial render. This may not be needed if we
    // have a suspense boundary BELOW the creation of the query client
    if (!browserQueryClient) {
      browserQueryClient = makeQueryClient();
    }

    return browserQueryClient;
  }
}
