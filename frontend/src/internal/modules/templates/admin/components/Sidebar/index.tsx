'use client';

import {
  drawerClasses,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Drawer as MuiDrawer,
  styled,
  Link as MuiLink,
} from '@mui/material';
import NextLink from 'next/link';
import { usePathname } from 'next/navigation';
import { PropsWithChildren, ReactNode } from 'react';

interface SidebarProps {
  items: {
    label: string;
    href: string;
    icon: ReactNode;
  }[];
}

const SidebarRoot = styled(MuiDrawer)(({ theme }) => ({
  '--Sidebar-width': '240px',
  width: 'var(--Sidebar-width)',
  flexShrink: 0,
  zIndex: `calc(${theme.vars.zIndex.appBar} - 1)`,

  [`& .${drawerClasses.paper}`]: {
    width: 'var(--Sidebar-width)',
  },
}));

const SidebarContent = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  flexGrow: 1,
  marginTop: `calc(8 * ${theme.vars.spacing})`,
  padding: `calc(2 * ${theme.vars.spacing})`,
}));

export function Sidebar(props: PropsWithChildren<SidebarProps>) {
  const pathname = usePathname();

  return (
    <SidebarRoot variant="permanent">
      <SidebarContent>
        <List dense disablePadding>
          {props.items.map((item) => (
            <MuiLink
              key={item.href}
              color="inherit"
              component={NextLink}
              href={item.href}
              underline="none"
            >
              <ListItem disablePadding sx={{ pb: 0.5 }}>
                <ListItemButton selected={pathname.startsWith(item.href)}>
                  <ListItemIcon>{item.icon}</ListItemIcon>
                  <ListItemText
                    primary={item.label}
                    primaryTypographyProps={{ noWrap: true }}
                  />
                </ListItemButton>
              </ListItem>
            </MuiLink>
          ))}
        </List>
      </SidebarContent>
    </SidebarRoot>
  );
}
