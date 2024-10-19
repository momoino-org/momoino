/**
 * Name of the cookie used to store the identity of a logged-in user.
 */
export const IdentityCookie = 'MOMOINO_IDENTITY';

/**
 * Name of the cookie used to manage the user's session data.
 */
export const SessionCookie = 'MOMOINO_SESSION';

/**
 * Name of the cookie used to store temporary login session information, typically used for login flow tracking.
 */
export const LoginSessionCookie = 'MOMOINO_LOGIN_SESSION';

/**
 * Name of the cookie used to store the CSRF token, providing protection against CSRF attacks.
 */
export const CsrfCookie = 'MOMOINO_CSRF';

/**
 * The header name for the CSRF token, required for CSRF validation in requests.
 */
export const CsrfHeaderName = 'X-Csrf-Token';
