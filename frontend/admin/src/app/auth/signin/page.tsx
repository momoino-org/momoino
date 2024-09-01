'use client';

import { valibotResolver } from '@hookform/resolvers/valibot';
import { Button, Card, Stack, TextField, Typography } from '@mui/material';
import { Controller, useForm } from 'react-hook-form';
import * as v from 'valibot';
import { useRouter, useSearchParams } from 'next/navigation';
import { useLoginByCredentials } from '@/internal/modules/core/ui/providers/AppProvider';
import { useTranslations } from 'next-intl';

const LoginSchema = v.object({
  email: v.pipe(v.string(), v.email()),
  password: v.pipe(v.string(), v.minLength(8)),
});

type LoginData = v.InferOutput<typeof LoginSchema>;

export default function LoginPage() {
  const t = useTranslations();
  const router = useRouter();
  const searchParams = useSearchParams();
  const { mutateAsync: loginAsync, status: loginStatus } =
    useLoginByCredentials();

  const { control, handleSubmit } = useForm<LoginData>({
    defaultValues: {
      email: '',
      password: '',
    },
    resolver: valibotResolver(LoginSchema),
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
            name="email"
            render={({ field: { ref, ...field }, fieldState }) => (
              <TextField
                autoComplete="email"
                error={fieldState.invalid}
                helperText={fieldState.error?.message}
                id="email"
                inputRef={ref}
                label={t('common.email')}
                type="email"
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
      </Card>
    </Stack>
  );
}
