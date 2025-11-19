/** @type {import('next').NextConfig} */
const nextConfig = {
  /* config options here */
  allowedDevOrigins: [
    "reflectively-credited-lorikeet.cloudpub.ru",
    "*.cloudpub.ru",
  ],
  reactCompiler: true,
  basePath: process.env.NEXT_PUBLIC_BASE_PATH || "",
  reactStrictMode: false,
}

module.exports = nextConfig
