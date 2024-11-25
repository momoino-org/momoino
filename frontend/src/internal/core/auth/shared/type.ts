import { z } from 'zod';
import { ProfileSchema } from './schema';

export type Profile = z.infer<typeof ProfileSchema>;
