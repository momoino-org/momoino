import { z } from 'zod';
import { Options } from 'ky';
import {
  http,
  PaginatedResponse,
  PaginatedResponseSchema,
  SingleResponse,
} from '@/internal/core/http';

export enum Provider {
  Google = 'google',
}

const ProviderDTOSchema = z
  .object({
    id: z.string(),
    name: z.string(),
    isEnabled: z.boolean(),
    createdAt: z.coerce.date(),
    createdBy: z.string(),
  })
  .strict();

type ProviderDTO = z.infer<typeof ProviderDTOSchema>;

export const CreateProviderParamsSchema = z
  .object({
    provider: z.string().min(1),
    clientId: z.string().min(1),
    clientSecret: z.string().min(1),
    redirectUrl: z.string().min(1),
    scopes: z.array(z.string()).min(1),
    isEnabled: z.boolean(),
  })
  .strict();

export type CreateProviderParams = z.infer<typeof CreateProviderParamsSchema>;

export async function getProviders(
  options?: Options,
): Promise<PaginatedResponse<ProviderDTO>> {
  const response = await http.get('api/v1/providers', options).json();
  return PaginatedResponseSchema(ProviderDTOSchema).parseAsync(response);
}

export async function createProvider(params: CreateProviderParams) {
  const response = await http
    .post('api/v1/providers', {
      json: params,
    })
    .json();

  return SingleResponse(ProviderDTOSchema).parseAsync(response);
}
