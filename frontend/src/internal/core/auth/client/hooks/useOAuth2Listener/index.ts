import { useEffect } from 'react';
import { z } from 'zod';
import { toast } from '@/internal/core/ui';

const AuthenticationMessageDataSchema = z
  .object({
    source: z.literal('useOAuth2'),
    payload: z.discriminatedUnion('status', [
      z.object({ status: z.literal('success') }),
      z.object({ status: z.literal('error'), details: z.unknown() }),
    ]),
  })
  .strict();

export type AuthenticationMessageData = z.infer<
  typeof AuthenticationMessageDataSchema
>;

export function useOAuth2Listener() {
  useEffect(() => {
    const authenticationTabListener = (event: MessageEvent) => {
      const data = AuthenticationMessageDataSchema.safeParse(event.data);

      if (!data.success) {
        return;
      }

      sessionStorage.removeItem('auth.provider');
      sessionStorage.removeItem('auth.state');
      sessionStorage.removeItem('auth.usePkce');
      sessionStorage.removeItem('auth.verifier');

      if (data.data.payload.status === 'error') {
        console.error(
          'Failed to authenticate. Please try again.',
          data.data.payload.details,
        );
        toast({
          severity: data.data.payload.status,
          message: 'Failed to authenticate. Please try again.',
        });
      } else {
        window.location.reload();
      }
    };

    window.addEventListener('message', authenticationTabListener);

    return () => {
      window.removeEventListener('message', authenticationTabListener);
    };
  }, []);
}
