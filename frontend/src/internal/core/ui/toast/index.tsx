'use client';

import { Alert, AlertColor } from '@mui/material';
import { toast as sonner } from 'sonner';

interface ToastOptions {
  severity: AlertColor;
  message: string;
}

export { Toaster } from 'sonner';

/**
 * Utility for managing toast notifications. Provides methods to display custom toast messages
 * and dismiss active toasts.
 */
export const notification = {
  /**
   * Displays a toast notification with a custom alert component.
   *
   * @param options - An object of type `ToastOptions` that specifies the severity level and message for the toast.
   * @returns The result of the `sonner.custom` call, which renders the custom toast.
   *
   * @example
   * ```typescript
   * notification.toast({ severity: 'success', message: 'Operation successful!' });
   * ```
   */
  toast(options: Readonly<ToastOptions>): ReturnType<typeof sonner.custom> {
    return sonner.custom((t) => (
      <Alert
        severity={options.severity}
        variant="filled"
        onClose={() => sonner.dismiss(t)}
      >
        {options.message}
      </Alert>
    ));
  },

  /**
   * Dismisses an active toast notification by its ID.
   * If no ID is provided, it dismisses all active toasts.
   *
   * @param id - An optional identifier for the specific toast to dismiss. If omitted, all toasts are dismissed.
   * @returns The result of the `sonner.dismiss` call, which dismisses the specified toast(s).
   *
   * @example
   * ```typescript
   * notification.dismiss(); // Dismisses all toasts
   * notification.dismiss(1); // Dismisses toast with ID 1
   * ```
   */
  dismiss(id?: string | number): ReturnType<typeof sonner.dismiss> {
    return sonner.dismiss(id);
  },
} as const;
