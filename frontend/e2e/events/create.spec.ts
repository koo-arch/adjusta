import { expect, test } from '../fixtures/auth';

// FullCalendarを複数ページで同時に初期化するとCI環境のメモリを圧迫するため、この画面だけ直列で実行する。
test.describe.configure({ mode: 'serial' });

test('[EVENT-009] イベント作成フォームの基本情報を表示できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/new');

    await expect(page.getByRole('heading', { name: 'イベント作成' })).toBeVisible();
    await expect(page.getByRole('heading', { name: '基本情報' })).toBeVisible();
    await expect(page.getByRole('textbox', { name: /タイトル/ })).toBeVisible();
    await expect(page.getByRole('textbox', { name: '場所' })).toBeVisible();
    await expect(page.getByRole('textbox', { name: '説明' })).toBeVisible();
    await expect(page.getByRole('button', { name: '次へ' })).toBeVisible();
});

test('[EVENT-010] 基本情報から候補日程ステップへ移動できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/new');
    await page.getByRole('button', { name: '次へ' }).click();

    await expect(page.getByRole('heading', { name: '候補日程' })).toBeVisible();
    await expect(
        page.getByText('カレンダーで範囲を選択するか、「日時を追加」から登録できます'),
    ).toBeVisible();
    await expect(page.getByRole('button', { name: '登録する' })).toBeVisible();
});

test('[EVENT-011] ステップを移動しても入力した基本情報を保持する', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/new');
    await page.getByRole('textbox', { name: /タイトル/ }).fill('E2E定例会議');
    await page.getByRole('textbox', { name: '場所' }).fill('オンライン');
    await page.getByRole('textbox', { name: '説明' }).fill('入力保持の確認');

    await page.getByRole('button', { name: '次へ' }).click();
    await expect(page.getByText('E2E定例会議')).toBeVisible();
    await expect(page.getByText('オンライン')).toBeVisible();
    await page.getByRole('button', { name: '戻る' }).click();

    await expect(page.getByRole('textbox', { name: /タイトル/ })).toHaveValue('E2E定例会議');
    await expect(page.getByRole('textbox', { name: '場所' })).toHaveValue('オンライン');
    await expect(page.getByRole('textbox', { name: '説明' })).toHaveValue('入力保持の確認');
});

test('[EVENT-012] タイトルと候補日程が未入力の場合はvalidationを表示する', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/new');
    await page.getByRole('button', { name: '次へ' }).click();
    await page.getByRole('button', { name: '登録する' }).click();

    await expect(page.getByText('日程は1つ以上選択してください')).toBeVisible();
    await page.getByRole('button', { name: '基本情報に入力エラーがあります' }).click();
    await expect(page.getByText('タイトルは必須です')).toBeVisible();
    await expect(page.getByRole('textbox', { name: /タイトル/ })).toHaveAttribute(
        'aria-invalid',
        'true',
    );
});

test('[EVENT-013] 候補日程を追加してイベントを作成できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events/new');
    await page.getByRole('textbox', { name: /タイトル/ }).fill('E2E作成イベント');
    await page.getByRole('textbox', { name: '場所' }).fill('オンライン');
    await page.getByRole('textbox', { name: '説明' }).fill('作成成功の確認');
    await page.getByRole('button', { name: '次へ' }).click();

    await page.getByRole('button', { name: '日時を追加' }).click();
    await expect(page.getByLabel('第1候補')).toBeVisible();

    const createRequestPromise = page.waitForRequest(
        (request) =>
            request.method() === 'POST' && request.url().endsWith('/api/calendar/event/draft'),
    );
    await page.getByRole('button', { name: '登録する' }).click();

    const createRequest = await createRequestPromise;
    const payload = createRequest.postDataJSON();
    expect(payload).toMatchObject({
        form_type: 'draft',
        title: 'E2E作成イベント',
        location: 'オンライン',
        description: '作成成功の確認',
    });
    expect(payload.selected_dates).toHaveLength(1);
    expect(payload.selected_dates[0]).toMatchObject({ id: null, priority: 1 });

    await expect(page).toHaveURL('/events/created-event', { timeout: 15_000 });
    await expect(page.getByRole('heading', { name: 'E2E作成イベント' })).toBeVisible();
});
