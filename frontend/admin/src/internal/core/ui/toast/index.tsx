'use client';

import { Alert, AlertColor } from '@mui/material';
import { toast as sonner } from 'sonner';

interface Toast {
  severity: AlertColor;
  message: string;
}

export { Toaster } from 'sonner';

export function toast(options: Toast): void {
  sonner.custom((t) => (
    <Alert
      severity={options.severity}
      variant="filled"
      onClose={() => sonner.dismiss(t)}
    >
      {options.message}
    </Alert>
  ));
}
