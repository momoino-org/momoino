import { getUserProfile } from '@/internal/modules/auth/services';
import { AccountMenu, TemplateFrame } from '@/internal/modules/core/ui';
import { MenuRounded, ModeNightRounded } from '@mui/icons-material';
import { AppBar, Toolbar, IconButton, Stack, Box } from '@mui/material';
import { PropsWithChildren } from 'react';

export default async function AdminLayout(props: PropsWithChildren) {
  const userProfile = await getUserProfile();

  if (userProfile === null) {
    return null;
  }

  return (
    <TemplateFrame>
      <AppBar elevation={1}>
        <Toolbar
          sx={{
            display: 'flex',
            justifyContent: 'space-between',
            width: '100%',
          }}
        >
          <IconButton aria-label="Open menu">
            <MenuRounded />
          </IconButton>
          <Stack direction="row" spacing={1}>
            <IconButton aria-label="Switch to the dark mode">
              <ModeNightRounded />
            </IconButton>
            <AccountMenu
              slotProps={{
                avatar: {
                  src: `https://ui-avatars.com/api/?rounded=true&name=${userProfile.username}&size=24`,
                },
                tooltip: {
                  title: 'Account',
                },
              }}
            />
          </Stack>
        </Toolbar>
      </AppBar>
      <Box sx={{ flex: '1 1', overflow: 'auto' }}>{props.children}</Box>
    </TemplateFrame>
  );
}
