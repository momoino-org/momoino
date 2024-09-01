import { createTheme, CssVarsThemeOptions, ThemeOptions } from '@mui/material';
import { Inter } from 'next/font/google';
import { colors } from './colors';

const inter = Inter({
  subsets: ['latin'],
  display: 'swap',
});

type DesignTokens = Omit<ThemeOptions, 'components'> &
  Pick<
    CssVarsThemeOptions,
    'defaultColorScheme' | 'colorSchemes' | 'components'
  > & {
    cssVariables?:
      | boolean
      | Pick<
          CssVarsThemeOptions,
          | 'colorSchemeSelector'
          | 'disableCssColorScheme'
          | 'cssVarPrefix'
          | 'shouldSkipGeneratingVar'
        >;
  };

const customTheme = createTheme();

export const brand = {
  50: 'hsl(220, 100%, 100%)',
  100: 'hsl(220, 100%, 92%)',
  200: 'hsl(220, 100%, 84%)',
  300: 'hsl(222, 100%, 76%)',
  400: 'hsl(224, 100%, 70%)',
  500: 'hsl(225, 100%, 60%)',
  600: 'hsl(226, 72%, 50%)',
  700: 'hsl(228, 76%, 41%)',
  800: 'hsl(230, 80%, 32%)',
  900: 'hsl(231, 86%, 26%)',
};

export const getDesignTokens = (): DesignTokens => ({
  shape: {
    borderRadius: 8,
  },
  cssVariables: {
    cssVarPrefix: '',
    colorSchemeSelector: 'class',
  },
  colorSchemes: {
    light: {
      palette: {
        primary: {
          light: brand[400],
          main: brand[600],
          dark: brand[800],
          contrastText: brand[50],
        },
        info: {
          light: '#D7F3FE',
          main: '#88D1FE',
          dark: '#2A7BD8',
          contrastText: colors.gray[50],
        },
      },
    },
    dark: {
      palette: {
        primary: {
          light: brand[300],
          main: brand[500],
          dark: brand[700],
          contrastText: brand[50],
        },
      },
    },
  },
  typography: {
    fontFamily: inter.style.fontFamily,
    fontSize: 14,
    h1: {
      fontSize: customTheme.typography.pxToRem(48),
      fontWeight: 600,
      lineHeight: 1.2,
      letterSpacing: -0.5,
    },
    h2: {
      fontSize: customTheme.typography.pxToRem(36),
      fontWeight: 600,
      lineHeight: 1.2,
    },
    h3: {
      fontSize: customTheme.typography.pxToRem(30),
      lineHeight: 1.2,
    },
    h4: {
      fontSize: customTheme.typography.pxToRem(24),
      fontWeight: 600,
      lineHeight: 1.5,
    },
    h5: {
      fontSize: customTheme.typography.pxToRem(20),
      fontWeight: 600,
    },
    h6: {
      fontSize: customTheme.typography.pxToRem(18),
      fontWeight: 600,
    },
    subtitle1: {
      fontSize: customTheme.typography.pxToRem(18),
    },
    subtitle2: {
      fontSize: customTheme.typography.pxToRem(14),
      fontWeight: 500,
    },
    body1: {
      fontSize: customTheme.typography.pxToRem(14),
    },
    body2: {
      fontSize: customTheme.typography.pxToRem(14),
      fontWeight: 400,
    },
    caption: {
      fontSize: customTheme.typography.pxToRem(12),
      fontWeight: 400,
    },
  },
});
