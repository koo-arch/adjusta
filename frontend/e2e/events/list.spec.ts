import { expect, test } from '../fixtures/auth';

test('[EVENT-001] イベントがない場合は作成案内を表示する', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events');

    await expect(page.getByRole('heading', { name: 'イベント一覧' })).toBeVisible();
    await expect(page.getByText('イベントがまだありません')).toBeVisible();
    await expect(page.getByText('候補日程を登録して、日程調整を始めましょう。')).toBeVisible();
});

test('[EVENT-002] イベント一覧から作成画面へ移動できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events');
    await expect(page.getByText('イベントがまだありません')).toBeVisible();
    await page.getByRole('link', { name: 'イベントを作成' }).last().click();

    await expect(page).toHaveURL('/events/new', { timeout: 15_000 });
    await expect(page.getByRole('heading', { name: 'イベント作成' })).toBeVisible();
});

test('[EVENT-003] タイトルで検索して該当なしの状態を表示できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events');
    await page.getByRole('textbox', { name: 'タイトルで検索' }).fill('定例会議');
    await page.getByRole('textbox', { name: 'タイトルで検索' }).press('Enter');

    await expect(page).toHaveURL('/events?title=%E5%AE%9A%E4%BE%8B%E4%BC%9A%E8%AD%B0');
    await expect(page.getByText('「定例会議」に一致するイベントはありません。')).toBeVisible();
});

test('[EVENT-004] ステータスで絞り込んで該当なしの状態を表示できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events');
    await page.getByRole('tab', { name: '調整中' }).click();

    await expect(page).toHaveURL('/events?status=active');
    await expect(page.getByText('「調整中」のイベントはありません。')).toBeVisible();
});
