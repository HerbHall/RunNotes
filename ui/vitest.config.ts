import path from "node:path";
import { defineConfig, mergeConfig } from "vitest/config";
import viteConfig from "./vite.config";

export default mergeConfig(
  viteConfig,
  defineConfig({
    resolve: {
      alias: {
        "@docker/extension-api-client": path.resolve(
          __dirname,
          "src/__mocks__/@docker/extension-api-client.ts",
        ),
      },
    },
    test: {
      environment: "jsdom",
      globals: true,
      setupFiles: "./src/test-setup.ts",
    },
  }),
);
