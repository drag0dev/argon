import { defineConfig, loadEnv } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';

const env = loadEnv();  // This will load all environment variables

export default defineConfig({
  plugins: [pluginReact()],
  source: {
    // Define global constants for the application, making them available as global variables
    define: {
      'process.env.API_URL': JSON.stringify(process.env.API_URL),
      'process.env.APP_POOL_ID': JSON.stringify(process.env.APP_POOL_ID),
      'process.env.POOL_CLIENT_ID': JSON.stringify(process.env.POOL_CLIENT_ID),
    }
  }
});

