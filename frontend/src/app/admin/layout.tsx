import {
  HomeRounded,
  MenuRounded,
  ModeNightRounded,
  MovieRounded,
  TvRounded,
} from '@mui/icons-material';
import {
  AppBar,
  Toolbar,
  IconButton,
  Stack,
  Box,
  Container,
} from '@mui/material';
import { PropsWithChildren } from 'react';
import { Sidebar, TemplateFrame } from '@/internal/modules/templates/admin';
import { AXccountMenu } from '@/internal/modules/templates/admin/providers/AccountMenu';

export default async function AdminLayout(props: PropsWithChildren) {
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
            <AXccountMenu />
          </Stack>
        </Toolbar>
      </AppBar>
      <Box sx={{ flex: '1 1', overflow: 'auto' }}>
        <Box sx={{ display: 'flex', height: `calc(100dvh - 64px)` }}>
          <Sidebar
            items={[
              { label: 'Home', href: '/admin', icon: <HomeRounded /> },
              {
                label: 'Movies',
                href: '/admin/movies',
                icon: <MovieRounded />,
              },
              {
                label: 'TV Shows',
                href: '/admin/tv-shows',
                icon: <TvRounded />,
              },
              {
                label: 'User Management',
                href: '/admin/usermgt',
                icon: <TvRounded />,
              },
            ]}
          />
          <Container maxWidth="xl" sx={{ flexGrow: 1, py: 2, height: '100%' }}>
            {props.children}
          </Container>
        </Box>
      </Box>
    </TemplateFrame>
  );
}
