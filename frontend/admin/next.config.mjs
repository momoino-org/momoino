import createNextIntlPlugin from 'next-intl/plugin';
import bundleAnalyzer from '@next/bundle-analyzer';

const withNextIntl = createNextIntlPlugin(
  './src/internal/core/i18n/request.ts',
);

// eslint-disable-next-line jsdoc/check-tag-names
/** @type {import('next').NextConfig} */
const nextConfig = {
  output: undefined,
  cleanDistDir: true,
  excludeDefaultMomentLocales: true,
  poweredByHeader: false,
};

// eslint-disable-next-line @typescript-eslint/no-var-requires
const withBundleAnalyzer = bundleAnalyzer({
  enabled: process.env.ANALYZE === 'true',
});

export default withBundleAnalyzer(withNextIntl(nextConfig));
