import type { NextConfig } from "next";
import createNextIntlPlugin from "next-intl/plugin";

const nextConfig: NextConfig = {
  images: {
    // Load images from randomuser.me for Lego driver profile pictures
    domains: ["randomuser.me"],
  },
  reactStrictMode: false,
};

const withNextIntl = createNextIntlPlugin();
export default withNextIntl(nextConfig);
