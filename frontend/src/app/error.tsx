'use client';

import { Alert, Button, Container, Stack, Typography } from '@mui/material';
import { useTranslations } from 'next-intl';

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  const t = useTranslations();

  const copyError = async () => {
    await navigator.clipboard.writeText(
      `// Name: ${error.name}\n\n// Message:\n${error.message}\n\n// Stack trace:\n${error.stack}\n\n// Cause:\n${error.cause ?? 'N/A'}\n\n// Digest:\n${error.digest ?? 'N/A'}`,
    );
  };

  return (
    <Stack style={{ height: '100dvh', justifyContent: 'center' }}>
      <Container>
        <Typography color="error" component="h1" variant="h4">
          {t('common.messages.internalError')}
        </Typography>
        <Alert icon={false} severity="error" variant="filled">
          Name: {error.name}
        </Alert>
        <Alert component="pre" icon={false} severity="error" variant="filled">
          {error.message}
        </Alert>
        <Alert component="pre" icon={false} severity="error" variant="filled">
          {error.stack}
        </Alert>
        <Stack direction="row-reverse" gap={1}>
          <Button color="error" variant="contained" onClick={() => copyError()}>
            Copy error
          </Button>

          <Button variant="outlined" onClick={() => reset()}>
            Try again
          </Button>
        </Stack>
      </Container>
    </Stack>
  );
}
