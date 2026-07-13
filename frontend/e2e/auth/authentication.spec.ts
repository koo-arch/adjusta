import { expect, test } from '@playwright/test';

const protectedRoutes = [
    { id: 'AUTH-001', route: '/dashboard' },
    { id: 'AUTH-002', route: '/events' },
    { id: 'AUTH-003', route: '/events/new' },
    { id: 'AUTH-004', route: '/events/test-event' },
    { id: 'AUTH-005', route: '/events/test-event/edit' },
    { id: 'AUTH-006', route: '/account' },
];

test.describe('未認証時の導線', () => {
    for (const { id, route } of protectedRoutes) {
        test(`[${id}] ${route} はログインページへ遷移する`, async ({ page }) => {
            await page.goto(route);

            await expect(page).toHaveURL('/login');
            await expect(page.getByRole('heading', { name: 'Adjusta' })).toBeVisible();
        });
    }
});

test('[AUTH-007] 期限切れセッションはcookieを削除してログインページへ遷移する', async ({
    context,
    page,
}) => {
    await context.addCookies([
        {
            name: 'session',
            value: 'expired-session',
            domain: 'localhost',
            path: '/',
            httpOnly: true,
            sameSite: 'Lax',
        },
    ]);

    await page.goto('/dashboard');

    await expect(page).toHaveURL('/login');
    await expect(page.getByRole('heading', { name: 'Adjusta' })).toBeVisible();
    await expect
        .poll(async () => (await context.cookies()).some((cookie) => cookie.name === 'session'))
        .toBe(false);
});
