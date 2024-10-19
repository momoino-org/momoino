import InitColorSchemeScript from '@mui/material/InitColorSchemeScript';
import { AppRouterCacheProvider } from '@mui/material-nextjs/v14-appRouter';
import { NextIntlClientProvider } from 'next-intl';
import { getLocale, getMessages } from 'next-intl/server';
import { Inter } from 'next/font/google';
import { PropsWithChildren } from 'react';
import { headers } from 'next/headers';
import { Metadata } from 'next';
import { SessionProvider } from 'next-auth/react';
import { MUIThemeProvider, StoreProvider, Toaster } from '@/internal/core/ui';
import { QueryClientProvider } from '@/internal/core/react-query';
import { config } from '@/internal/core/i18n/request';
import { auth, signIn } from '@/auth';

const interFont = Inter({
  weight: ['300', '400', '500', '700'],
  subsets: ['latin'],
  display: 'swap',
  variable: '--font-family',
});

export const fetchCache = 'default-cache';

export async function generateMetadata(): Promise<Metadata> {
  const csrfToken = headers().get('X-Csrf-Token')!;

  return {
    other: {
      'csrf-token': csrfToken,
    },
  };
}

export default async function RootLayout({ children }: PropsWithChildren) {
  const locale = await getLocale();
  const messages = await getMessages();
  const session = await auth();

  if (session === null) {
    return signIn();
  }

  return (
    <html suppressHydrationWarning dir="ltr" lang={locale}>
      <body className={interFont.variable}>
        <InitColorSchemeScript attribute="class" />

        <NextIntlClientProvider formats={config.formats} messages={messages}>
          <AppRouterCacheProvider options={{ enableCssLayer: true }}>
            <MUIThemeProvider locale={locale}>
              <Toaster
                pauseWhenPageIsHidden
                richColors
                closeButton={false}
                duration={6_000}
                position="top-center"
                toastOptions={{
                  unstyled: true,
                }}
              />

              <QueryClientProvider>
                <SessionProvider refetchInterval={30} session={session} >
                  <StoreProvider>{children}
                    {children}
                  </StoreProvider>
                </SessionProvider>

              </QueryClientProvider>
            </MUIThemeProvider>
          </AppRouterCacheProvider>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
