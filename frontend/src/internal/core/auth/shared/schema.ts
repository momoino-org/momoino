import { z } from 'zod';

export const JWTSchema = z.strictObject({
  accessToken: z.string(),
  refreshToken: z.string(),
});

export const ProfileSchema = z.strictObject({
  id: z.string(),
  username: z.string(),
  email: z.string(),
  firstName: z.string(),
  lastName: z.string(),
  locale: z.string(),
});
