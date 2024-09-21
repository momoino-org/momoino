'use client';

import { PropsWithChildren } from 'react';
import { CssBaseline, ThemeProvider } from '@mui/material';
import {
  enUS as coreEnUS,
  viVN as coreViVN,
  type Localization as MuiLocalization,
} from '@mui/material/locale';
import {
  enUS as dataGridEnUS,
  viVN as dataGridViVN,
} from '@mui/x-data-grid/locales';
import { type Localization as MuiDataGridLocalization } from '@mui/x-data-grid/internals';
import { createTheme } from './theme';

interface MuiThemeProvider {
  locale: string;
}

const translations: Record<
  string,
  (MuiLocalization | MuiDataGridLocalization)[]
> = {
  en: [dataGridEnUS, coreEnUS],
  vi: [dataGridViVN, coreViVN],
};

export function MUIThemeProvider(props: PropsWithChildren<MuiThemeProvider>) {
  return (
    <ThemeProvider theme={createTheme(translations[props.locale])}>
      <CssBaseline />
      {props.children}
    </ThemeProvider>
  );
}
