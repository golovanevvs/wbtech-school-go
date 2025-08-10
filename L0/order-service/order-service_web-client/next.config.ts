import type { NextConfig } from "next"

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:6000/:path*",
      },
      {
        source: "/order/:path*",
        destination: "http://localhost:6000/order/:path*",
      },
    ]
  },
}

module.exports = nextConfig
