import { z } from 'zod';
import { JWTSchema, ProfileSchema } from './schema';

export type JWT = z.infer<typeof JWTSchema>;

export type Profile = z.infer<typeof ProfileSchema>;

export interface JwtPayload {
  sid: string;
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
