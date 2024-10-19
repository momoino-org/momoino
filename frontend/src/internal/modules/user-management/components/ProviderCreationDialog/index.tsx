import { zodResolver } from '@hookform/resolvers/zod';
import { Controller, SubmitHandler, useForm } from 'react-hook-form';
import {
  DialogTitle,
  DialogContent,
  TextField,
  DialogActions,
  Button,
  Autocomplete,
  FormControlLabel,
  Checkbox,
  MenuItem,
} from '@mui/material';
import { useMutation } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import {
  createProvider,
  CreateProviderParams,
  CreateProviderParamsSchema,
  Provider,
} from '../../service';
import { EmptyResponse } from '@/internal/core/http';
import { toast, useModalContext } from '@/internal/core/ui';

export function ProviderCreationDialog() {
  const t = useTranslations();
  const modal = useModalContext();

  const { control, handleSubmit } = useForm<CreateProviderParams>({
    resolver: zodResolver(CreateProviderParamsSchema),
    defaultValues: {
      provider: '',
      clientId: '',
      clientSecret: '',
      redirectUrl: '',
      scopes: [],
      isEnabled: false,
    },
  });

  const { mutateAsync, isPending } = useMutation({
    mutationKey: ['create-a-provider'],
    mutationFn: createProvider,
    onSuccess: async (result) => {
      const response = await EmptyResponse.strip().parseAsync(result);

      toast({
        severity: 'success',
        message: response.message,
      });

      modal.close();
    },
  });

  const onSubmit: SubmitHandler<CreateProviderParams> = async (data) => {
    await mutateAsync(data);
  };

  return (
    <>
      <DialogTitle>Create a provider</DialogTitle>
      <DialogContent
        dividers
        sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}
      >
        <Controller
          control={control}
          name="provider"
          render={({ field: { ref, ...field }, fieldState }) => (
            <TextField
              select
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
              id="provider"
              inputRef={ref}
              label="Provider"
              variant="filled"
              {...field}
            >
              <MenuItem value={Provider.Google}>Google</MenuItem>
            </TextField>
          )}
        />

        <Controller
          control={control}
          name="clientId"
          render={({ field: { ref, ...field }, fieldState }) => (
            <TextField
              autoComplete="new-password"
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
              id="clientId"
              inputRef={ref}
              label="Client ID"
              type="password"
              variant="filled"
              {...field}
            />
          )}
        />

        <Controller
          control={control}
          name="clientSecret"
          render={({ field: { ref, ...field }, fieldState }) => (
            <TextField
              autoComplete="new-password"
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
              id="clientSecret"
              inputRef={ref}
              label="Client Secret"
              type="password"
              variant="filled"
              {...field}
            />
          )}
        />

        <Controller
          control={control}
          name="redirectUrl"
          render={({ field: { ref, ...field }, fieldState }) => (
            <TextField
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
              id="redirectUrl"
              inputRef={ref}
              label="Redirect URL"
              variant="filled"
              {...field}
            />
          )}
        />

        <Controller
          control={control}
          name="scopes"
          render={({ field, fieldState }) => (
            <Autocomplete
              freeSolo
              multiple
              id="scopes"
              options={[]}
              renderInput={(params) => (
                <TextField
                  {...params}
                  error={fieldState.invalid}
                  helperText={fieldState.error?.message}
                  label="Scopes"
                  variant="filled"
                />
              )}
              onChange={(e, data) => field.onChange(data)}
            />
          )}
        />

        <Controller
          control={control}
          name="isEnabled"
          render={({ field: { onChange, value } }) => (
            <FormControlLabel
              control={<Checkbox checked={value} onChange={onChange} />}
              label="Is enabled?"
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
          onClick={handleSubmit(onSubmit, (data) => {
            console.log(data);
          })}
        >
          {t('common.create')}
        </Button>
      </DialogActions>
    </>
  );
}
