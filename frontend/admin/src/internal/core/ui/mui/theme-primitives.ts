import {
  alpha,
  createTheme,
  CssVarsThemeOptions,
  ThemeOptions,
} from '@mui/material';
import { colors } from './colors';

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
          light: colors.brand[400],
          main: colors.brand[600],
          dark: colors.brand[800],
          contrastText: colors.brand[50],
        },
        info: {
          light: '#D7F3FE',
          main: '#88D1FE',
          dark: '#2A7BD8',
          contrastText: colors.gray[50],
        },
        divider: alpha(colors.gray[400], 0.4),
        background: {
          default: 'hsl(0, 0%, 99%)',
          paper: 'hsl(220, 35%, 97%)',
        },
      },
    },
    dark: {
      palette: {
        primary: {
          light: colors.brand[300],
          main: colors.brand[500],
          dark: colors.brand[700],
          contrastText: colors.brand[50],
        },
      },
    },
  },
  typography: {
    fontFamily: 'var(--font-family)',
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
      fontWeight: 600,
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
