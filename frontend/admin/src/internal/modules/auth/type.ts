export interface JwtPayload {
  sub: string;
  exp: number;
  nbf: number;
  iat: number;
  email: string;
  given_name: string;
  family_name: string;
  preferred_username: string;
  locale: string;
  roles: string[];
  permissions: string[];
}
