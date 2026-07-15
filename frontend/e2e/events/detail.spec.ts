import { expect, test } from '../fixtures/auth';

// 編集画面のFullCalendar初期化と詳細画面を並列に開くとCI環境のメモリを圧迫するため、直列で実行する。
test.describe.configure({ mode: 'serial' });

test('[EVENT-014] イベント詳細と候補日程を表示できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/detail-event');

    await expect(page.getByRole('heading', { name: '詳細確認イベント' })).toBeVisible();
    await expect(page.getByText('会議室A')).toBeVisible();
    await expect(page.getByText('イベント詳細の説明')).toBeVisible();
    await expect(page.getByRole('heading', { name: '候補日程' })).toBeVisible();
    await expect(page.getByLabel('第1候補').locator('..')).toContainText('7月22日(水) 1:00 - 2:00');
    await expect(page.getByLabel('第2候補').locator('..')).toContainText('7月21日(火) 1:00 - 2:00');
    await expect(page.getByLabel('第3候補').locator('..')).toContainText('7月20日(月) 1:00 - 2:00');
    await expect(page.getByText('未同期', { exact: true }).first()).toBeVisible();
});

test('[EVENT-015] 存在しないイベントでは404表示になる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/missing-event');

    await expect(page.getByText('イベントが見つかりませんでした。')).toBeVisible({
        timeout: 15_000,
    });
    await expect(page.getByRole('link', { name: 'イベント一覧へ戻る' })).toHaveAttribute(
        'href',
        '/events',
    );
});

test('[EVENT-016] イベント詳細から編集画面へ移動できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/detail-event');
    await expect(page.getByRole('heading', { name: '詳細確認イベント' })).toBeVisible();
    await page.getByRole('link', { name: '編集' }).click();

    await expect(page).toHaveURL('/events/detail-event/edit', { timeout: 15_000 });
    await expect(page.getByRole('heading', { name: 'イベント編集' })).toBeVisible();
});
