import { defineConfig, devices } from '@playwright/test';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const checkoutStoreDir = path.resolve(__dirname, '..');
const repoDir = path.resolve(checkoutStoreDir, '..', '..');
const muaraDataDir = path.join(checkoutStoreDir, '.muara-e2e');

export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: 'list',
  use: {
    baseURL: 'http://127.0.0.1:8080',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: [
    {
      command: `mkdir -p ${muaraDataDir}/data && MUARA_SERVER_PORT=9001 go run ${repoDir}/cmd/muara start --config ${checkoutStoreDir}/e2e/muara.config.yml`,
      url: 'http://127.0.0.1:9001/healthz',
      timeout: 60_000,
      reuseExistingServer: !process.env.CI,
    },
    {
      command: `mailpit -l "127.0.0.1:9035" -s "127.0.0.1:9025" --smtp-auth-accept-any --smtp-auth-allow-insecure`,
      url: 'http://127.0.0.1:9035/api/v1/info',
      timeout: 30_000,
      reuseExistingServer: !process.env.CI,
    },
    {
      command: `OPENMUARA_URL=http://127.0.0.1:9001 MAILPIT_PORT=9025 go run .`,
      cwd: checkoutStoreDir,
      url: 'http://127.0.0.1:8080/api/product',
      timeout: 60_000,
      reuseExistingServer: !process.env.CI,
    },
  ],
});
