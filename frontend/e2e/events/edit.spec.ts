import { expect, test } from '../fixtures/auth';

// FullCalendarを複数ページで同時に初期化するとCI環境のメモリを圧迫するため、この画面だけ直列で実行する。
test.describe.configure({ mode: 'serial' });

test('[EVENT-018] イベント編集フォームに現在の内容を表示できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/edit-event/edit');

    await expect(page.getByRole('heading', { name: 'イベント編集' })).toBeVisible();
    await expect(page.getByRole('textbox', { name: /タイトル/ })).toHaveValue('編集前イベント');
    await expect(page.getByRole('textbox', { name: '場所' })).toHaveValue('会議室B');
    await expect(page.getByRole('textbox', { name: '説明' })).toHaveValue('編集前の説明');

    await page.getByRole('button', { name: '次へ' }).click();
    await expect(page.getByLabel('第1候補')).toBeVisible();
});

test('[EVENT-019] 基本情報を更新して詳細画面へ戻れる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/edit-event/edit');
    await page.getByRole('textbox', { name: /タイトル/ }).fill('編集後イベント');
    await page.getByRole('textbox', { name: '場所' }).fill('オンライン');
    await page.getByRole('textbox', { name: '説明' }).fill('編集後の説明');
    await page.getByRole('button', { name: '次へ' }).click();

    const updateRequestPromise = page.waitForRequest(
        (request) =>
            request.method() === 'PUT' &&
            request.url().endsWith('/api/calendar/event/draft/edit-event'),
    );
    await page.getByRole('button', { name: '保存する' }).click();

    const updateRequest = await updateRequestPromise;
    expect(updateRequest.postDataJSON()).toMatchObject({
        id: 'edit-event',
        form_type: 'edit',
        title: '編集後イベント',
        location: 'オンライン',
        description: '編集後の説明',
        status: 'active',
        proposed_dates: [
            {
                id: '11111111-1111-4111-8111-111111111111',
                priority: 1,
            },
        ],
    });
    await expect(page).toHaveURL('/events/edit-event', { timeout: 15_000 });
});

test('[EVENT-020] 編集画面で候補日程を削除して追加できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/edit-event/edit');
    await page.getByRole('button', { name: '次へ' }).click();
    await page.getByRole('button', { name: 'この日程を削除', exact: true }).click();

    await expect(
        page.getByText('カレンダーで範囲を選択するか、「日時を追加」から登録できます'),
    ).toBeVisible();

    await page.getByRole('button', { name: '日時を追加' }).click();
    await page.getByRole('button', { name: '編集を完了' }).click();
    await expect(page.getByLabel('第1候補')).toBeVisible();
});

test('[EVENT-021] 存在しないイベントの編集画面では404表示になる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/missing-event/edit');

    await expect(page.getByText('イベントが見つかりませんでした。')).toBeVisible({
        timeout: 15_000,
    });
    await expect(page.getByRole('link', { name: 'イベント一覧へ戻る' })).toHaveAttribute(
        'href',
        '/events',
    );
});
