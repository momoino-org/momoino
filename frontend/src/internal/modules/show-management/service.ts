import { z } from 'zod';
import { type Options } from 'ky';
import {
  http,
  PaginatedResponseSchema,
  PaginatedResponse,
  SingleResponse,
} from '@/internal/core/http';

export enum ShowKind {
  Movie = 'movie',
  TVShow = 'tv_show',
}

const ShowDTOSchema = z
  .object({
    id: z.string(),
    kind: z.nativeEnum(ShowKind),
    originalTitle: z.string().min(1).max(255),
    originalOverview: z.string().min(1).max(255),
    originalLanguage: z.string().min(2).max(2),
    keywords: z.array(z.string()),
    isReleased: z.boolean(),
    createdAt: z.coerce.date(),
    updatedAt: z.coerce.date(),
  })
  .strict();

type ShowDTO = z.infer<typeof ShowDTOSchema>;

export const CreateShowFormDataSchema = z
  .object({
    kind: z.nativeEnum(ShowKind),
    originalTitle: z.string().min(1).max(255),
    originalOverview: z.string().min(1).max(255),
    originalLanguage: z.string().min(2).max(2),
    keywords: z.array(z.string()),
    isReleased: z.boolean(),
  })
  .strict();

export type CreateShowFormData = z.infer<typeof CreateShowFormDataSchema>;

/**
 * Fetches a paginated list of shows from the API.
 *
 * @param options - Additional options to pass to the HTTP request.
 * @returns A promise that resolves to a paginated list of show data.
 */
export async function getShows(
  options: Options = {},
): Promise<PaginatedResponse<ShowDTO>> {
  const response = await http.get('api/v1/shows', options).json();
  return PaginatedResponseSchema(ShowDTOSchema).parseAsync(response);
}

/**
 * Creates a new show in the API
 *
 * @param data - The data for the new show.
 * @returns A promise that resolves to the created show data.
 */
export async function createShow(
  data: CreateShowFormData,
): Promise<SingleResponse<ShowDTO>> {
  const response = await http
    .post('api/v1/shows', {
      json: data,
    })
    .json();

  return SingleResponse(ShowDTOSchema).parseAsync(response);
}
