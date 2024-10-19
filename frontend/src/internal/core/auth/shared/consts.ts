/**
 * Name of the cookie used to store the CSRF token, providing protection against CSRF attacks.
 */
export const CsrfCookie = 'MOMOINO_CSRF';

/**
 * The header name for the CSRF token, required for CSRF validation in requests.
 */
export const CsrfHeaderName = 'X-Csrf-Token';
