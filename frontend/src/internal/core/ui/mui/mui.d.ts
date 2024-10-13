import type {} from '@mui/x-data-grid';
import { ButtonProps } from '@mui/material';

declare module '@mui/x-data-grid' {
  interface ToolbarPropsOverrides {
    toolbarRefreshButton?: ButtonProps;
  }
}
