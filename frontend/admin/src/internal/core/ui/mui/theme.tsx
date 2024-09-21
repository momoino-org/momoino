'use client';

import type {} from '@mui/material/themeCssVarsAugmentation';
import type {} from '@mui/x-data-grid/themeAugmentation';
import {
  alpha,
  checkboxClasses,
  createTheme as muiCreateTheme,
  responsiveFontSizes,
  tablePaginationClasses,
  touchRippleClasses,
} from '@mui/material';
import { gridClasses } from '@mui/x-data-grid';
import {
  CheckBoxOutlineBlankRounded,
  CheckRounded,
  RemoveRounded,
} from '@mui/icons-material';
import { type Localization as MuiLocalization } from '@mui/material/locale';
import { type Localization as MuiDataGridLocalization } from '@mui/x-data-grid/internals';
import { colors } from './colors';
import { getDesignTokens } from './theme-primitives';

export const createTheme = (
  locale: (MuiLocalization | MuiDataGridLocalization)[],
) =>
  responsiveFontSizes(
    muiCreateTheme(
      {
        ...getDesignTokens(),
        components: {
          MuiToolbar: {
            styleOverrides: {
              root: ({ theme }) => ({
                [theme.breakpoints.up('sm')]: {
                  paddingLeft: `calc(2 * ${theme.vars.spacing})`,
                  paddingRight: `calc(2 * ${theme.vars.spacing})`,
                },
              }),
            },
          },
          MuiListItemIcon: {
            styleOverrides: {
              root: ({ theme }) => ({
                color: theme.vars.palette.grey[800],
                minWidth: 0,
              }),
            },
          },
          MuiListItemButton: {
            styleOverrides: {
              root: ({ theme }) => ({
                gap: theme.vars.spacing,
                borderRadius: theme.vars.shape.borderRadius,
              }),
            },
          },
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
                [`& .${touchRippleClasses.root} .${touchRippleClasses.child}`]:
                  {
                    borderRadius: theme.vars.shape.borderRadius,
                  },
              }),
            },
          },
          MuiCheckbox: {
            defaultProps: {
              icon: (
                <CheckBoxOutlineBlankRounded
                  sx={{ color: 'hsla(210, 0%, 0%, 0.0)' }}
                />
              ),
              checkedIcon: <CheckRounded sx={{ height: 14, width: 14 }} />,
              indeterminateIcon: (
                <RemoveRounded sx={{ height: 14, width: 14 }} />
              ),
            },
            styleOverrides: {
              root: ({ theme }) => ({
                margin: 10,
                height: 16,
                width: 16,
                borderRadius: 5,
                border: '1px solid ',
                borderColor: alpha(colors.gray[300], 0.8),
                boxShadow: '0 0 0 1.5px hsla(210, 0%, 0%, 0.04) inset',
                backgroundColor: alpha(colors.gray[100], 0.4),
                transition: 'border-color, background-color, 120ms ease-in',
                '&:hover': {
                  borderColor: theme.vars.palette.primary.main,
                },
                '&.Mui-focusVisible': {
                  outline: `3px solid ${alpha(colors.brand[500], 0.5)}`,
                  outlineOffset: '2px',
                  borderColor: colors.brand[400],
                },
                '&.Mui-checked': {
                  color: 'white',
                  backgroundColor: colors.brand[500],
                  borderColor: colors.brand[500],
                  boxShadow: `none`,
                  '&:hover': {
                    backgroundColor: colors.brand[600],
                  },
                },
                ...theme.applyStyles('dark', {
                  borderColor: alpha(colors.gray[700], 0.8),
                  boxShadow: '0 0 0 1.5px hsl(210, 0%, 0%) inset',
                  backgroundColor: alpha(colors.gray[900], 0.8),
                  '&:hover': {
                    borderColor: colors.brand[300],
                  },
                  '&.Mui-focusVisible': {
                    borderColor: colors.brand[400],
                    outline: `3px solid ${alpha(colors.brand[500], 0.5)}`,
                    outlineOffset: '2px',
                  },
                }),
              }),
            },
          },
          MuiDataGrid: {
            // TODO: https://github.com/mui/mui-x/issues/14708
            // defaultProps: {
            //   slotProps: {
            //     loadingOverlay: {
            //       variant: 'linear-progress',
            //       noRowsVariant: 'linear-progress',
            //     },
            //   },
            // },
            styleOverrides: {
              root: ({ theme }) => ({
                '--DataGrid-overlayHeight': '300px',
                borderColor: theme.vars.palette.divider,
                backgroundColor: theme.vars.palette.background.default,
                [`& .${gridClasses.columnHeader}`]: {
                  backgroundColor: theme.vars.palette.background.paper,
                },
                [`& .${gridClasses.footerContainer}`]: {
                  backgroundColor: theme.vars.palette.background.paper,
                },
                [`& .${checkboxClasses.root}`]: {
                  padding: theme.spacing(0.5),
                },
                [`& .${tablePaginationClasses.root}`]: {
                  marginRight: theme.spacing(1),
                  '& .MuiIconButton-root': {
                    maxHeight: 32,
                    maxWidth: 32,
                    '& > svg': {
                      fontSize: '1rem',
                    },
                  },
                },
              }),
              cell: ({ theme }) => ({
                borderTopColor: theme.vars.palette.divider,
              }),
              row: ({ theme }) => ({
                '&:last-of-type': {
                  borderBottom: `1px solid ${(theme.vars || theme).palette.divider}`,
                },
                '&:hover': {
                  backgroundColor: (theme.vars || theme).palette.action.hover,
                },
                '&.Mui-selected': {
                  background: (theme.vars || theme).palette.action.selected,
                  '&:hover': {
                    backgroundColor: (theme.vars || theme).palette.action.hover,
                  },
                },
              }),
              columnsManagementHeader: ({ theme }) => ({
                paddingRight: theme.spacing(3),
                paddingLeft: theme.spacing(3),
              }),
              columnHeaderTitleContainer: {
                flexGrow: 1,
                justifyContent: 'space-between',
              },
              columnHeaderDraggableContainer: { paddingRight: 2 },
            },
          },
        },
      },
      // Adding a translation here causes performance issues.
      // It might be related to this issue https://github.com/mui/mui-x/issues/14708
      ...locale,
    ),
  );
