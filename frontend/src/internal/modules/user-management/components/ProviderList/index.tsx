'use client';

import { GridPaginationModel, GridColDef } from '@mui/x-data-grid';
import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { useFormatter } from 'next-intl';
import { isEqual } from 'radash';
import { useMemo } from 'react';
import { useDebounceValue } from 'usehooks-ts';
import { Button } from '@mui/material';
import { getProviders } from '../../service';
import { ProviderCreationDialog } from '../ProviderCreationDialog';
import { DataGrid } from '@/internal/core/ui/mui/components/DataGrid';
import { useModal } from '@/internal/core/ui';

export function ProviderList() {
  const formatter = useFormatter();
  const creationDialog = useModal({
    content: <ProviderCreationDialog />,
    onClose: () => {
      refetch();
    },
  });

  const [paginationModel, setPaginationModel] =
    useDebounceValue<GridPaginationModel>(
      {
        page: 0,
        pageSize: 10,
      },
      500,
      { equalityFn: isEqual },
    );

  const { data, isFetching, refetch } = useQuery({
    queryKey: ['get-shows', paginationModel.page, paginationModel.pageSize],
    placeholderData: keepPreviousData,
    queryFn: ({ signal }) =>
      getProviders({
        signal,
        searchParams: {
          page: paginationModel.page + 1,
          pageSize: paginationModel.pageSize,
        },
      }),
  });

  const handlePaginationModelChange = (
    paginationModel: GridPaginationModel,
  ) => {
    setPaginationModel(paginationModel);
  };

  const columns = useMemo<GridColDef[]>(
    () => [
      {
        field: 'id',
        type: 'string',
        headerName: 'ID',
        flex: 1,
      },
      {
        field: 'name',
        type: 'string',
        headerName: 'Name',
        flex: 1,
        hideable: false,
      },
      {
        field: 'isEnabled',
        type: 'boolean',
        headerName: 'Enabled',
        flex: 1,
        headerAlign: 'left',
        align: 'left',
      },
      {
        field: 'createdAt',
        type: 'dateTime',
        headerName: 'Created At',
        flex: 1,
        valueFormatter: (value) => formatter.dateTime(value, 'short'),
      },
      {
        field: 'createdBy',
        type: 'string',
        headerName: 'Created By',
        flex: 1,
      },
    ],
    [formatter],
  );

  return (
    <>
      <Button onClick={creationDialog.open}>Create</Button>
      <DataGrid
        checkboxSelection
        disableColumnResize
        columns={columns}
        getRowId={(row) => row.id}
        loading={isFetching}
        pageSizeOptions={[10, 25, 50, 100]}
        paginationMode="server"
        paginationModel={paginationModel}
        rowCount={data?.pagination.totalRows ?? -1}
        rows={data?.data}
        slotProps={{
          toolbar: {
            toolbarRefreshButton: {
              onClick: () => refetch(),
            },
          },
        }}
        onPaginationModelChange={handlePaginationModelChange}
      />
      {creationDialog.content}
    </>
  );
}
