'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { Button, Card, Stack, TextField, Typography } from '@mui/material';
import { Controller, useForm } from 'react-hook-form';
import { useRouter, useSearchParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { z } from 'zod';
import { useMutation } from '@tanstack/react-query';
import { Google } from '@mui/icons-material';
import { loginByCredentials, OAuth2Button } from '@/internal/core/auth/client';
import { backendOrigin } from '@/internal/core/config';

const LoginSchema = z.strictObject({
  username: z.string(),
  password: z.string().min(8),
});

type LoginData = z.infer<typeof LoginSchema>;

export default function LoginPage() {
  const t = useTranslations();
  const router = useRouter();
  const searchParams = useSearchParams();
  const { mutateAsync: loginAsync, status: loginStatus } = useMutation({
    mutationKey: ['login-by-credentials'],
    mutationFn: (params: LoginData) =>
      loginByCredentials(params.username, params.password),
  });

  const { control, handleSubmit } = useForm<LoginData>({
    defaultValues: {
      username: '',
      password: '',
    },
    resolver: zodResolver(LoginSchema),
  });

  const onSubmit = async (data: LoginData) => {
    await loginAsync(data, {
      onSuccess: () => {
        router.push(searchParams.get('redirectTo') ?? '/');
      },
    });
  };

  return (
    <Stack
      spacing={1}
      sx={(theme) => ({
        height: '100dvh',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        backgroundImage:
          'radial-gradient(ellipse at 50% 50%, hsl(210, 100%, 97%), hsl(0, 0%, 100%))',
        ...theme.applyStyles('dark', {
          backgroundImage:
            'radial-gradient(at 50% 50%, hsla(210, 100%, 16%, 0.5), hsl(220, 30%, 5%))',
        }),
      })}
    >
      <Card
        sx={(theme) => ({
          display: 'flex',
          flexDirection: 'column',
          alignSelf: 'center',
          width: '100%',
          padding: theme.spacing(4),
          gap: theme.spacing(2),
          margin: 'auto',
          [theme.breakpoints.up('sm')]: {
            maxWidth: '450px',
          },
          boxShadow:
            'hsla(220, 30%, 5%, 0.05) 0px 5px 15px 0px, hsla(220, 25%, 10%, 0.05) 0px 15px 35px -5px',
          ...theme.applyStyles('dark', {
            boxShadow:
              'hsla(220, 30%, 5%, 0.5) 0px 5px 15px 0px, hsla(220, 25%, 10%, 0.08) 0px 15px 35px -5px',
          }),
        })}
        variant="outlined"
      >
        <Typography component="h1" variant="h4">
          {t('page.signin.title')}
        </Typography>
        <Stack
          noValidate
          component="form"
          spacing={2}
          onSubmit={handleSubmit(onSubmit)}
        >
          <Controller
            control={control}
            name="username"
            render={({ field: { ref, ...field }, fieldState }) => (
              <TextField
                autoComplete="username"
                error={fieldState.invalid}
                helperText={fieldState.error?.message}
                id="username"
                inputRef={ref}
                label={t('common.username')}
                type="text"
                variant="filled"
                {...field}
              />
            )}
          />

          <Controller
            control={control}
            name="password"
            render={({ field: { ref, ...field }, fieldState }) => (
              <TextField
                autoComplete="current-password"
                error={fieldState.invalid}
                helperText={fieldState.error?.message}
                id="password"
                inputRef={ref}
                label={t('common.password')}
                type="password"
                variant="filled"
                {...field}
              />
            )}
          />

          <Button
            fullWidth
            disabled={loginStatus === 'pending'}
            type="submit"
            variant="contained"
          >
            {t('page.signin.signinBtn')}
          </Button>
        </Stack>

        <OAuth2Button
          pkce
          authURL={new URL('api/v1/login/providers/google', backendOrigin)}
          provider="google"
        >
          <Google />
        </OAuth2Button>
      </Card>
    </Stack>
  );
}
