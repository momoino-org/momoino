'use client';

import NextLink from 'next/link';
import {
  Box,
  Breadcrumbs,
  Button,
  Chip,
  Link,
  Stack,
  Typography,
} from '@mui/material';
import {
  DataGrid,
  GridActionsCellItem,
  GridColDef,
  GridPaginationModel,
} from '@mui/x-data-grid';
import { keepPreviousData, useQuery } from '@tanstack/react-query';
import { CreateRounded } from '@mui/icons-material';
import { useFormatter, useTranslations } from 'next-intl';
import { useMemo } from 'react';
import { useDebounceValue } from 'usehooks-ts';
import { isEqual } from 'radash';
import { useModal } from '@/internal/core/ui';
import {
  CreateShowDialog,
  getShows,
  ShowKind,
} from '@/internal/modules/show-management';

export default function MovieListPage() {
  const t = useTranslations();
  const formatter = useFormatter();
  const creationDialog = useModal({
    content: <CreateShowDialog />,
    onClose: () => {
      refetch();
    },
  });
  const [paginationModel, setPaginationModel] =
    useDebounceValue<GridPaginationModel>(
      {
        page: 0,
        pageSize: 100,
      },
      500,
      { equalityFn: isEqual },
    );

  const { data, isFetching, refetch } = useQuery({
    queryKey: ['get-shows', paginationModel.page, paginationModel.pageSize],
    placeholderData: keepPreviousData,
    queryFn: ({ signal }) =>
      getShows({
        signal,
        searchParams: {
          page: paginationModel.page + 1,
          pageSize: paginationModel.pageSize,
        },
      }),
  });

  const columns: GridColDef[] = useMemo(
    () => [
      {
        field: 'kind',
        type: 'singleSelect',
        headerName: 'Kind',
        flex: 1,
        valueOptions: [
          {
            value: ShowKind.Movie,
            label: t('page.movies.creationDialog.fields.kind.options.movie'),
          },
          {
            value: ShowKind.TVShow,
            label: t('page.movies.creationDialog.fields.kind.options.tv_show'),
          },
        ],
        renderCell: (params) => {
          return (
            <Chip
              label={params.value}
              size="small"
              sx={{ textTransform: 'uppercase' }}
            />
          );
        },
      },
      {
        field: 'originalTitle',
        headerName: 'Original title',
        flex: 1,
      },
      {
        field: 'originalLanguage',
        headerName: 'Original language',
        flex: 1,
      },
      {
        field: 'isReleased',
        type: 'boolean',
        headerName: 'Is released',
        headerAlign: 'left',
        align: 'left',
        flex: 1,
      },
      {
        field: 'createdAt',
        type: 'dateTime',
        headerName: 'Created at',
        flex: 1,
        valueFormatter: (value) => formatter.dateTime(value, 'short'),
      },
      {
        field: 'updatedAt',
        type: 'dateTime',
        headerName: 'Updated at',
        flex: 1,
        valueFormatter: (value) => formatter.dateTime(value, 'short'),
      },
      {
        field: 'actions',
        type: 'actions',
        headerName: 'Actions',
        width: 80,
        getActions: (params) => [
          <GridActionsCellItem
            key={`${params.id}-details`}
            icon={<CreateRounded />}
            label="Details"
            onClick={() => {
              params.id;
            }}
          />,
        ],
      },
    ],
    [t, formatter],
  );

  const handlePaginationModelChange = (
    paginationModel: GridPaginationModel,
  ) => {
    setPaginationModel(paginationModel);
  };

  return (
    <>
      <Stack height="100%" spacing={1}>
        <Breadcrumbs>
          <Link
            color="inherit"
            component={NextLink}
            href="/admin"
            underline="hover"
          >
            Home
          </Link>
          <Link
            color="inherit"
            component={NextLink}
            href="/admin/movies"
            underline="hover"
          >
            Movies
          </Link>
        </Breadcrumbs>
        <Typography component="h1" variant="h3">
          Movies
        </Typography>
        <Stack direction="row-reverse" spacing={1}>
          <Button
            disabled={creationDialog.visible}
            variant="contained"
            onClick={creationDialog.open}
          >
            Create
          </Button>

          <Button
            disabled={isFetching}
            variant="contained"
            onClick={() => refetch()}
          >
            Refresh
          </Button>
        </Stack>
        <Box sx={{ height: '600px' }}>
          <DataGrid
            checkboxSelection
            disableColumnResize
            columns={columns}
            getRowId={(row) => row.id}
            loading={isFetching}
            pageSizeOptions={[10, 25, 50, 100]}
            paginationMode="server"
            paginationModel={{
              page: data?.pagination.page ? data.pagination.page - 1 : 0,
              pageSize: data?.pagination.pageSize
                ? data.pagination.pageSize
                : 10,
            }}
            rowCount={data?.pagination.totalRows ?? -1}
            rows={data?.data}
            onPaginationModelChange={handlePaginationModelChange}
          />
        </Box>
      </Stack>
    </>
  );
}
