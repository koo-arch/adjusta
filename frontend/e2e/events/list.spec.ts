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
    await page.getByRole('main').getByRole('link', { name: '新規作成' }).first().click();

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

test('[EVENT-005] タイトル検索条件をクリアできる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events?title=%E5%AE%9A%E4%BE%8B%E4%BC%9A%E8%AD%B0');
    await expect(page.getByText('「定例会議」に一致するイベントはありません。')).toBeVisible();

    await page.getByRole('button', { name: '検索条件をクリア' }).last().click();

    await expect(page).toHaveURL('/events');
    await expect(page.getByText('イベントがまだありません')).toBeVisible();
});

test('[EVENT-006] ステータスを保持したままタイトル検索できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events?status=active');
    const searchInput = page.getByRole('textbox', { name: 'タイトルで検索' });
    await searchInput.fill('定例会議');
    await searchInput.press('Enter');

    await expect(page).toHaveURL((url) => {
        return (
            url.pathname === '/events' &&
            url.searchParams.get('status') === 'active' &&
            url.searchParams.get('title') === '定例会議'
        );
    });
    await expect(page.getByText('「定例会議」に一致するイベントはありません。')).toBeVisible();
});

test('[EVENT-007] 検索結果のイベントカードを表示できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events?title=%E8%A1%A8%E7%A4%BA%E7%A2%BA%E8%AA%8D');

    const eventLink = page.getByRole('link', { name: /表示確認イベント/ });
    await expect(eventLink).toBeVisible();
    await expect(eventLink).toHaveAttribute('href', '/events/visible-event');
    await expect(eventLink.getByText('調整中', { exact: true })).toBeVisible();
    await expect(eventLink.getByText('候補日程はありません')).toBeVisible();
});

test('[EVENT-008] 検索結果の次ページへ移動できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
}) => {
    await page.goto('/events?title=%E3%83%9A%E3%83%BC%E3%82%B8%E3%83%B3%E3%82%B0');
    await expect(page.getByText('1ページ目のイベント')).toBeVisible();

    await page.getByRole('link', { name: '次のページへ' }).click();

    await expect(page).toHaveURL((url) => {
        return url.searchParams.get('title') === 'ページング' && url.searchParams.get('page') === '2';
    });
    await expect(page.getByText('2ページ目のイベント')).toBeVisible();
    await expect(page.getByText('21–21 / 21 件')).toBeVisible();
});
