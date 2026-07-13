import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
    testDir: './e2e',
    fullyParallel: true,
    forbidOnly: Boolean(process.env.CI),
    retries: process.env.CI ? 2 : 0,
    workers: process.env.CI ? 1 : undefined,
    reporter: 'html',
    use: {
        baseURL: 'http://localhost:3100',
        trace: 'on-first-retry',
    },
    projects: [
        {
            name: 'chromium',
            use: { ...devices['Desktop Chrome'] },
        },
    ],
    webServer: [
        {
            command: 'node e2e/fixtures/mock-backend.mjs',
            url: 'http://localhost:3101/health',
            reuseExistingServer: !process.env.CI,
            timeout: 30_000,
        },
        {
            command: 'INTERNAL_BACKEND_URL=http://localhost:3101 NEXT_DIST_DIR=.next-e2e npm run dev -- --port 3100',
            url: 'http://localhost:3100',
            reuseExistingServer: !process.env.CI,
            timeout: 120_000,
        },
    ],
});
