import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  experimental: {
    after: true,
  },
  images: {
    remotePatterns: [
      {
        hostname: 'localhost'
      }
    ]
  }
};

export default nextConfig;
