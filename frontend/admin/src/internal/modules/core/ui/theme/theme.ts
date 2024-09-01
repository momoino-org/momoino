'use client';

import type {} from '@mui/material/themeCssVarsAugmentation';
import {
  createTheme,
  responsiveFontSizes,
  touchRippleClasses,
} from '@mui/material';
import { getDesignTokens } from './theme-primitives';

export const theme = responsiveFontSizes(
  createTheme({
    ...getDesignTokens(),
    components: {
      MuiAppBar: {
        styleOverrides: {
          root: ({ theme }) => ({
            '--AppBar-background': theme.palette.background.paper,
            position: 'relative',
            flex: '0 0 auto',
          }),
        },
      },
      MuiFilledInput: {
        defaultProps: {
          disableUnderline: true,
        },
        styleOverrides: {
          input: {
            '&:-webkit-autofill': {
              borderRadius: 'inherit',
            },
          },
          root: ({ theme }) => ({
            borderRadius: theme.vars.shape.borderRadius,
          }),
        },
      },
      MuiIconButton: {
        styleOverrides: {
          root: ({ theme }) => ({
            borderRadius: theme.vars.shape.borderRadius,
            border: '1px solid',
            borderColor: theme.vars.palette.divider,
            [`& .${touchRippleClasses.root} .${touchRippleClasses.child}`]: {
              borderRadius: theme.vars.shape.borderRadius,
            },
          }),
        },
      },
    },
  }),
);
