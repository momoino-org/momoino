'use client';

import { Button, ButtonProps } from '@mui/material';
import {
  DataGrid as MuiDataGrid,
  DataGridProps,
  GridToolbarContainer,
  GridToolbarProps,
  GridToolbarColumnsButton,
  GridToolbarFilterButton,
  GridToolbarDensitySelector,
} from '@mui/x-data-grid';
import { forwardRef } from 'react';

function GridToolbarRefreshButton(props: ButtonProps) {
  return <Button {...props}>Refresh</Button>;
}

const GridToolbar = forwardRef<HTMLDivElement, GridToolbarProps>(
  function GridToolbar(props, ref) {
    const { toolbarRefreshButton, ...rest } = props;

    return (
      <GridToolbarContainer ref={ref} {...rest}>
        <GridToolbarColumnsButton />
        <GridToolbarFilterButton />
        <GridToolbarDensitySelector />
        {toolbarRefreshButton && (
          <GridToolbarRefreshButton onClick={toolbarRefreshButton.onClick} />
        )}
      </GridToolbarContainer>
    );
  },
);

export function DataGrid(props: DataGridProps) {
  return (
    <MuiDataGrid
      {...props}
      slots={{
        toolbar: GridToolbar,
        ...props.slots,
      }}
    />
  );
}
