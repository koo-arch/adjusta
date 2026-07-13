import { expect, test } from '../fixtures/auth';

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
            url: 'http://localhost:3100',
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

test('[AUTH-008] 画面表示後にセッションが切れた場合は次のページ要求でログインページへ遷移する', async ({
    authenticatedSession,
    context,
    page,
    request,
}) => {
    await page.goto('/account');
    await expect(page.getByRole('heading', { name: 'アカウント設定' })).toBeVisible();

    const expireResponse = await request.post(
        `http://localhost:3101/__e2e/sessions/${authenticatedSession.token}/expire`,
    );
    expect(expireResponse.ok()).toBe(true);

    await page.reload();

    await expect(page).toHaveURL('/login');
    await expect
        .poll(async () => (await context.cookies()).some((cookie) => cookie.name === 'session'))
        .toBe(false);
});

test('[AUTH-009] 画面表示後のAPI操作が401になった場合はログインページへ遷移する', async ({
    authenticatedSession,
    context,
    page,
    request,
}) => {
    await page.goto('/account');
    const candidateSync = page.getByRole('switch', {
        name: '候補日程の Google Calendar 同期',
    });
    await expect(candidateSync).toBeEnabled({ timeout: 15_000 });

    const expireResponse = await request.post(
        `http://localhost:3101/__e2e/sessions/${authenticatedSession.token}/expire`,
    );
    expect(expireResponse.ok()).toBe(true);

    const updateResponsePromise = page.waitForResponse(
        (response) =>
            response.request().method() === 'PUT' &&
            response.url().endsWith('/api/calendar-settings/candidate-sync'),
    );
    await candidateSync.click();

    const updateResponse = await updateResponsePromise;
    expect(updateResponse.status()).toBe(401);
    await expect(page).toHaveURL('/login');
    await expect
        .poll(async () => (await context.cookies()).some((cookie) => cookie.name === 'session'))
        .toBe(false);
});

for (const { id, route } of [
    { id: 'AUTH-010', route: '/' },
    { id: 'AUTH-011', route: '/login' },
]) {
    test(`[${id}] ログイン済みで ${route} を開くとダッシュボードへ遷移する`, async ({
        authenticatedSession: _authenticatedSession,
        page,
    }) => {
        await page.goto(route);

        await expect(page).toHaveURL('/dashboard');
        await expect(page.getByRole('heading', { name: 'ホーム' })).toBeVisible();
    });
}
