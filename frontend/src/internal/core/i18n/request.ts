import { getRequestConfig } from 'next-intl/server';
import { IntlConfig } from 'next-intl';
import { auth } from '@/auth';

export const config: Pick<IntlConfig, 'formats'> = {
  formats: {
    dateTime: {
      short: {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false,
      },
    },
  },
};

export default getRequestConfig(async () => {
  const session = await auth();
  const locale = session?.user?.locale ?? 'en';

  return {
    locale,
    messages: (await import(`./messages/${locale}.json`)).default,
    ...config,
  };
});
