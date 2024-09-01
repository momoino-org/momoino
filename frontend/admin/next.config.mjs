import createNextIntlPlugin from 'next-intl/plugin';

const withNextIntl = createNextIntlPlugin(
  './src/internal/modules/core/i18n/request.ts',
);

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: undefined,
  cleanDistDir: true,
  excludeDefaultMomentLocales: true,
  poweredByHeader: false,
};

export default withNextIntl(nextConfig);
