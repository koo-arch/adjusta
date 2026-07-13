import { expect, test } from '../fixtures/auth';

test('[DASHBOARD-001] 対応が必要なイベントがない状態を表示できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/dashboard');

    await expect(page.getByRole('heading', { name: 'ホーム' })).toBeVisible();
    await expect(page.getByText('対応が必要なイベントはありません')).toBeVisible();
});

test('[DASHBOARD-002] 直近の予定がない状態を表示できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/dashboard');
    await page.getByRole('tab', { name: '直近の予定' }).click();

    await expect(page.getByText('直近の予定はありません')).toBeVisible();
});
