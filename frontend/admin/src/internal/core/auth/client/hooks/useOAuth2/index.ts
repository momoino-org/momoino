import { useCallback } from 'react';

function sha256(message: string): Promise<ArrayBuffer> {
  const encoder = new TextEncoder();
  const data = encoder.encode(message);
  return crypto.subtle.digest('SHA-256', data);
}

function bytesToBase64(bytes: Uint8Array) {
  const binString = String.fromCodePoint(...bytes);
  return btoa(binString);
}

function urlEncodeB64(input: string): string {
  const b64Chars: Record<string, string> = { '+': '-', '/': '_', '=': '' };
  return input.replace(/[\+\/=]/g, (m) => b64Chars[m]);
}

function randomString(length: number): string {
  const validChars =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let array = new Uint8Array(length);
  crypto.getRandomValues(array);
  array = array.map((x) => validChars.charCodeAt(x % validChars.length));
  return String.fromCharCode(...array);
}

async function generatePKCE(length: number) {
  const codeVerifier = randomString(length);
  const codeChallenge = urlEncodeB64(
    bytesToBase64(new Uint8Array(await sha256(codeVerifier))),
  );

  return {
    codeVerifier,
    codeChallenge,
  };
}

export function useOAuth2(props: {
  provider: string;
  authURL: URL;
  usePkce: boolean;
}) {
  const handleClick = useCallback(async () => {
    const state = randomString(64);

    const url = new URL(props.authURL);
    url.searchParams.set('state', state);

    if (props.usePkce) {
      const pkce = await generatePKCE(64);
      url.searchParams.set('codeChallenge', pkce.codeChallenge);
      sessionStorage.setItem('auth.verifier', pkce.codeVerifier);
    } else {
      sessionStorage.removeItem('auth.verifier');
    }

    sessionStorage.setItem('auth.usePkce', String(props.usePkce));
    sessionStorage.setItem('auth.provider', props.provider);
    sessionStorage.setItem('auth.state', state);

    window.open(url);
  }, [props.authURL, props.provider, props.usePkce]);

  return {
    handleClick,
  };
}
