/**
 * Determines whether the current environment is a server-side environment.
 */
export const isServer: boolean = typeof window === 'undefined';

/**
 * Determines whether the current environment is a production environment.
 */
export const isProduction: boolean = process.env.NODE_ENV === 'production';

/**
 * Determines the appropriate backend URL based on the environment.
 *
 * If the code is running on the server (e.g., during server-side rendering),
 * it uses the `NEXT_BACKEND_HOST` environment variable.
 * Otherwise, on the client-side, it uses the `NEXT_PUBLIC_BACKEND_HOST` environment variable.
 */
export const backendOrigin: string = isServer
  ? process.env.NEXT_BACKEND_ORIGIN
  : process.env.NEXT_PUBLIC_BACKEND_ORIGIN;

/**
 * Retrieves the frontend origin URL based on the environment.
 */
export const frontendOrigin: string = process.env.NEXT_PUBLIC_FRONTEND_ORIGIN;
