import { Box } from '@mui/material';
import { PropsWithChildren, ReactNode } from 'react';

interface TemplateFrameProps extends PropsWithChildren {}

export function TemplateFrame(props: TemplateFrameProps): ReactNode {
  return (
    <Box
      sx={{
        height: '100dvh',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      {props.children}
    </Box>
  );
}
