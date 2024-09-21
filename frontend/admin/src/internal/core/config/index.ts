/**
 * Determines whether the current environment is a server-side environment.
 */
export const isServer: boolean = typeof window === 'undefined';

/**
 * Determines whether the current environment is a production environment.
 */
export const isProduction: boolean = process.env.NODE_ENV === 'production';
