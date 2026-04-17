import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:8080',
  },
  webServer: {
    command: 'go run ./webapp',
    url: 'http://localhost:8080/health',
    reuseExistingServer: false,
  },
});
