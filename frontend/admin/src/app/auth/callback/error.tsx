'use client';

import { Typography } from '@mui/material';

export default function ErrorBoundary() {
  return (
    <Typography>
      An error occurred while trying to authenticate. Please try again.
    </Typography>
  );
}
