import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Produces .next/standalone, a self-contained server the Docker image can run
  // without carrying node_modules.
  output: "standalone",
};

export default nextConfig;
