import { Button } from '@mui/material';
import { PropsWithChildren } from 'react';
import { useOAuth2 } from '../../hooks/useOAuth2';
import { useOAuth2Listener } from '../../hooks/useOAuth2Listener';

interface OAuth2ButtonProps {
  provider: string;
  authURL: URL;
  pkce?: boolean;
}

export function OAuth2Button(props: PropsWithChildren<OAuth2ButtonProps>) {
  const { handleClick } = useOAuth2({
    provider: props.provider,
    authURL: props.authURL,
    usePkce: props.pkce === true,
  });
  useOAuth2Listener();

  return <Button onClick={handleClick}>{props.children}</Button>;
}
