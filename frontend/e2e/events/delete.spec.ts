import { expect, test } from '../fixtures/auth';

test.describe.configure({ mode: 'serial' });

test('[EVENT-027] イベント削除の確認ダイアログをキャンセルできる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    const deleteRequests: string[] = [];
    page.on('request', (request) => {
        if (request.method() === 'DELETE') {
            deleteRequests.push(request.url());
        }
    });

    await page.goto('/events/delete-event');
    await page.getByRole('button', { name: '削除', exact: true }).click();

    const dialog = page.getByRole('alertdialog');
    await expect(dialog.getByRole('heading', { name: 'イベントを削除しますか?' })).toBeVisible();
    await expect(
        dialog.getByText('「削除対象イベント」を候補日程ごと削除します。この操作は取り消せません。'),
    ).toBeVisible();
    await dialog.getByRole('button', { name: 'キャンセル' }).click();

    await expect(dialog).toBeHidden();
    expect(deleteRequests).toHaveLength(0);
    await expect(page).toHaveURL('/events/delete-event');
});

test('[EVENT-028] イベントを削除して一覧画面へ移動できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/delete-event');
    await page.getByRole('button', { name: '削除', exact: true }).click();

    const deleteRequestPromise = page.waitForRequest(
        (request) =>
            request.method() === 'DELETE' &&
            request.url().endsWith('/api/calendar/event/draft/delete-event'),
    );
    await page.getByRole('alertdialog').getByRole('button', { name: '削除する' }).click();

    await deleteRequestPromise;
    await expect(page.getByText('イベントを削除しました')).toBeVisible();
    await expect(page).toHaveURL('/events', { timeout: 15_000 });
    await expect(page.getByRole('heading', { name: 'イベント一覧' })).toBeVisible();
});

test('[EVENT-029] 削除APIが失敗した場合はエラーを表示して詳細画面に留まる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/delete-error-event');
    await page.getByRole('button', { name: '削除', exact: true }).click();
    await page.getByRole('alertdialog').getByRole('button', { name: '削除する' }).click();

    await expect(page.getByText('イベントの削除処理に失敗しました')).toBeVisible();
    await expect(page).toHaveURL('/events/delete-error-event');
    await expect(page.getByRole('heading', { name: '削除失敗イベント' })).toBeVisible();
});
