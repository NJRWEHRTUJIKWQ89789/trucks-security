/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "cargomax-dashboard.vercel.app",
      },
    ],
  },
};

module.exports = nextConfig;
