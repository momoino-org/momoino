'use client';

import {
  Autocomplete,
  Button,
  Checkbox,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  MenuItem,
  TextField,
} from '@mui/material';
import { useMutation } from '@tanstack/react-query';
import { SubmitHandler, useForm, Controller } from 'react-hook-form';
import { useTranslations } from 'next-intl';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  CreateShowFormDataSchema,
  createShow,
  CreateShowFormData,
  ShowKind,
} from '../service';
import { EmptyResponse } from '@/internal/core/http';
import { notification, useModalContext } from '@/internal/core/ui';

export function CreateShowDialog() {
  const t = useTranslations();
  const modal = useModalContext();

  const { control, handleSubmit } = useForm<CreateShowFormData>({
    resolver: zodResolver(CreateShowFormDataSchema),
    defaultValues: {
      kind: ShowKind.Movie,
      originalTitle: '',
      originalOverview: '',
      originalLanguage: '',
      keywords: [],
      isReleased: false,
    },
  });

  const { mutateAsync, isPending } = useMutation({
    mutationKey: ['create-show'],
    mutationFn: createShow,
    onSuccess: async (result) => {
      const response = await EmptyResponse.strip().parseAsync(result);

      notification.toast({
        severity: 'success',
        message: response.message,
      });

      modal.close();
    },
  });

  const onSubmit: SubmitHandler<CreateShowFormData> = async (data) => {
    await mutateAsync(data);
  };

  return (
    <>
      <DialogTitle>{t('page.movies.creationDialog.title')}</DialogTitle>
      <DialogContent
        dividers
        sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}
      >
        <Controller
          control={control}
          name="kind"
          render={({ field: { ref, ...field }, fieldState }) => (
            <TextField
              select
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
              id="kind"
              inputRef={ref}
              label={t('page.movies.creationDialog.fields.kind.label')}
              variant="filled"
              {...field}
            >
              <MenuItem value={ShowKind.Movie}>
                {t('page.movies.creationDialog.fields.kind.options.movie')}
              </MenuItem>
              <MenuItem value={ShowKind.TVShow}>
                {t('page.movies.creationDialog.fields.kind.options.tv_show')}
              </MenuItem>
            </TextField>
          )}
        />

        <Controller
          control={control}
          name="originalTitle"
          render={({ field: { ref, ...field }, fieldState }) => (
            <TextField
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
              id="originalTitle"
              inputRef={ref}
              label={t('page.movies.creationDialog.fields.originalTitle')}
              variant="filled"
              {...field}
            />
          )}
        />

        <Controller
          control={control}
          name="originalOverview"
          render={({ field: { ref, ...field }, fieldState }) => (
            <TextField
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
              id="originalOverview"
              inputRef={ref}
              label={t('page.movies.creationDialog.fields.originalOverview')}
              variant="filled"
              {...field}
            />
          )}
        />

        <Controller
          control={control}
          name="originalLanguage"
          render={({ field: { ref, ...field }, fieldState }) => (
            <TextField
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
              id="originalLanguage"
              inputRef={ref}
              label={t('page.movies.creationDialog.fields.originalLanguage')}
              variant="filled"
              {...field}
            />
          )}
        />

        <Controller
          control={control}
          name="keywords"
          render={({ field }) => (
            <Autocomplete
              freeSolo
              multiple
              id="keywords"
              options={[]}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label={t('page.movies.creationDialog.fields.keywords')}
                  placeholder="Favorites"
                  variant="filled"
                />
              )}
              onChange={(e, data) => field.onChange(data)}
            />
          )}
        />

        <Controller
          control={control}
          name="isReleased"
          render={({ field: { onChange, value } }) => (
            <FormControlLabel
              control={<Checkbox checked={value} onChange={onChange} />}
              label={t('page.movies.creationDialog.fields.isReleased')}
            />
          )}
        />
      </DialogContent>
      <DialogActions>
        <Button disabled={isPending} variant="outlined" onClick={modal.close}>
          {t('common.close')}
        </Button>
        <Button
          disabled={isPending}
          type="submit"
          variant="contained"
          onClick={handleSubmit(onSubmit)}
        >
          {t('common.create')}
        </Button>
      </DialogActions>
    </>
  );
}
