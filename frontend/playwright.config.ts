import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
    testDir: './e2e',
    fullyParallel: true,
    forbidOnly: Boolean(process.env.CI),
    retries: process.env.CI ? 2 : 0,
    workers: process.env.CI ? 1 : undefined,
    reporter: [
        ['list'],
        ['html', { outputFolder: 'playwright-report', open: 'never' }],
        ['junit', { outputFile: 'test-results/e2e-results.xml' }],
    ],
    use: {
        baseURL: 'http://localhost:3100',
        screenshot: 'only-on-failure',
        trace: 'retain-on-failure',
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
            command: 'INTERNAL_BACKEND_URL=http://localhost:3101 NEXT_PUBLIC_API_BASE_URL= NEXT_DIST_DIR=.next-e2e npm run dev -- --port 3100',
            url: 'http://localhost:3100',
            reuseExistingServer: !process.env.CI,
            timeout: 120_000,
        },
    ],
});
