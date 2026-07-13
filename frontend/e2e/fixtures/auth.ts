import { expect, test as base } from '@playwright/test';

interface AuthFixtures {
    authenticatedSession: {
        token: string;
    };
}

export const test = base.extend<AuthFixtures>({
    authenticatedSession: async ({ context }, use, testInfo) => {
        const token = `active-session-${testInfo.workerIndex}-${Date.now()}`;
        await context.addCookies([
            {
                name: 'session',
                value: token,
                url: 'http://localhost:3100',
                httpOnly: true,
                sameSite: 'Lax',
            },
        ]);

        await use({ token });
    },
});

export { expect };
