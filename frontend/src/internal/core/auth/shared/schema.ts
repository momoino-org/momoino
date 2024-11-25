import { z } from 'zod';

export const ProfileSchema = z.strictObject({
  id: z.string(),
  username: z.string(),
  email: z.string(),
  emailVerified: z.boolean(),
  firstName: z.string(),
  lastName: z.string(),
  locale: z.string(),
});
