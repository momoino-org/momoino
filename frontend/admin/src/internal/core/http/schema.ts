import { z } from 'zod';

export const EmptyResponse = z
  .object({
    message: z.string(),
    messageId: z.string(),
    timestamp: z.string(),
  })
  .strict();

export type EmptyResponse = z.infer<typeof EmptyResponse>;

export function SingleResponse<T>(data: z.ZodType<T>) {
  return EmptyResponse.extend({
    data: data,
  }).strict();
}

export type SingleResponse<T> = z.infer<ReturnType<typeof SingleResponse<T>>>;

export function PaginatedResponseSchema<T>(data: z.ZodType<T>) {
  return SingleResponse(z.array(data))
    .extend({
      pagination: z.object({
        page: z.number(),
        pageSize: z.number(),
        totalRows: z.number(),
        totalPages: z.number(),
      }),
    })
    .strict();
}

export type PaginatedResponse<T> = z.infer<
  ReturnType<typeof PaginatedResponseSchema<T>>
>;
