import { getRequestConfig } from 'next-intl/server';
import { getUserProfile } from '../../auth/services';

export default getRequestConfig(async () => {
  const userProfileResponse = await getUserProfile();
  const locale = userProfileResponse?.locale ?? 'en';

  return {
    locale,
    messages: (await import(`./messages/${locale}.json`)).default,
  };
});
