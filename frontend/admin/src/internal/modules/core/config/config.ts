/**
 * Determines whether the application should use HTTPS based on the environment variable `USE_HTTPS`.
 *
 * @returns {boolean} Returns `true` if HTTPS should be used, otherwise `false`.
 */
export const isUseHttps: boolean = Boolean(process.env.USE_HTTPS);

/**
 * Determines whether the current environment is a server-side environment.
 *
 * @returns {boolean} Returns `true` if the current environment is a server-side environment, otherwise `false`.
 */
export const isServer: boolean = typeof window === 'undefined';
