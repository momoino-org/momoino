import { getRequestConfig } from 'next-intl/server';
import { IntlConfig } from 'next-intl';
import { getUserProfile } from '@/internal/core/auth/server';

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
  const userProfileResponse = await getUserProfile();
  const locale = userProfileResponse?.locale ?? 'en';

  return {
    locale,
    messages: (await import(`./messages/${locale}.json`)).default,
    ...config,
  };
});
