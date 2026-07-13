import { expect, test } from '../fixtures/auth';

test.describe.configure({ mode: 'serial' });

test('[EVENT-022] 日程確定ダイアログを表示できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/confirm-event');
    await page.getByRole('button', { name: '日程を確定' }).click();

    const dialog = page.getByRole('dialog');
    await expect(dialog.getByRole('heading', { name: '日程を確定' })).toBeVisible();
    await expect(dialog.getByRole('tab', { name: '候補から選択' })).toBeVisible();
    await expect(dialog.getByRole('tab', { name: '手動で入力' })).toBeVisible();
    await expect(dialog.getByRole('combobox')).toContainText('候補日程を選択');
});

test('[EVENT-023] 候補日程を未選択の場合はvalidationを表示する', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/confirm-event');
    await page.getByRole('button', { name: '日程を確定' }).click();
    await page.getByRole('dialog').getByRole('button', { name: '確定する' }).click();

    await expect(page.getByText('日程を選択してください', { exact: true })).toBeVisible();
    await expect(page.getByRole('dialog')).toBeVisible();
});

test('[EVENT-024] 候補日程を選択して確定できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/confirm-event');
    await page.getByRole('button', { name: '日程を確定' }).click();
    await page.getByRole('dialog').getByRole('combobox').click();
    await page.getByRole('option').first().click();

    const confirmRequestPromise = page.waitForRequest(
        (request) =>
            request.method() === 'PATCH' &&
            request.url().endsWith('/api/calendar/event/confirm/confirm-event'),
    );
    await page.getByRole('dialog').getByRole('button', { name: '確定する' }).click();

    const confirmRequest = await confirmRequestPromise;
    expect(confirmRequest.postDataJSON()).toEqual({
        confirm_date: {
            id: '22222222-2222-4222-8222-222222222222',
            google_event_id: 'candidate-google-event-1',
            start: '2026-07-22T01:00:00.000Z',
            end: '2026-07-22T02:00:00.000Z',
            priority: 1,
        },
    });

    await expect(page.getByText('日程を確定しました')).toBeVisible();
    await expect(page.getByRole('dialog')).toBeHidden();
    await expect(page.getByText('確定日時')).toBeVisible();
    await expect(page.getByRole('button', { name: '確定日程を変更' })).toBeVisible();
});

test('[EVENT-025] 確定APIが失敗した場合はエラーを表示してダイアログを維持する', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/confirm-error-event');
    await page.getByRole('button', { name: '日程を確定' }).click();
    await page.getByRole('dialog').getByRole('combobox').click();
    await page.getByRole('option').first().click();
    await page.getByRole('dialog').getByRole('button', { name: '確定する' }).click();

    await expect(page.getByText('日程の確定処理に失敗しました')).toBeVisible();
    await expect(page.getByRole('dialog')).toBeVisible();
});

test('[EVENT-026] 手動入力の日程が未入力の場合はvalidationを表示する', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/confirm-error-event');
    await page.getByRole('button', { name: '日程を確定' }).click();
    await page.getByRole('dialog').getByRole('tab', { name: '手動で入力' }).click();
    await page.getByRole('dialog').getByRole('button', { name: '確定する' }).click();

    await expect(page.getByText('開始日時は必須です')).toBeVisible();
    await expect(page.getByText('終了日時は必須です')).toBeVisible();
    await expect(page.getByRole('dialog')).toBeVisible();
});
