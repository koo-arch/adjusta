import { expect, test } from '../fixtures/auth';

// cacheComponents(PPR)の本番挙動を検証する。webServer が本番ビルドで
// 起動していることが前提(dev モードでは静的シェルが生成されない)。

const readUntil = async (
    reader: ReadableStreamDefaultReader<Uint8Array>,
    decoder: TextDecoder,
    predicate: (html: string) => boolean,
) => {
    let html = '';
    while (!predicate(html)) {
        const { done, value } = await reader.read();
        if (done) {
            break;
        }
        html += decoder.decode(value, { stream: true });
    }
    return html;
};

test('[PPR-001] 静的シェルがユーザー情報より先にストリーミングされる', async ({
    authenticatedSession,
    request,
}) => {
    const controlURL = `http://localhost:3101/__e2e/sessions/${authenticatedSession.token}/user`;
    const pauseResponse = await request.post(`${controlURL}/pause`);
    expect(pauseResponse.ok()).toBe(true);

    const controller = new AbortController();
    const decoder = new TextDecoder();
    let reader: ReadableStreamDefaultReader<Uint8Array> | undefined;

    try {
        const response = await fetch('http://localhost:3100/dashboard', {
            headers: { cookie: `session=${authenticatedSession.token}` },
            redirect: 'manual',
            signal: controller.signal,
        });
        expect(response.status).toBe(200);
        expect(response.body).not.toBeNull();
        reader = response.body!.getReader();

        const shellHTML = await readUntil(
            reader,
            decoder,
            (html) =>
                html.includes('イベント一覧') &&
                html.includes('ホーム') &&
                html.includes('ユーザー情報を読み込み中'),
        );

        // /users/me は保留中なので、最初のHTMLには静的シェルとfallbackだけが入る
        expect(shellHTML).toContain('イベント一覧');
        expect(shellHTML).toContain('ホーム');
        expect(shellHTML).toContain('ユーザー情報を読み込み中');
        expect(shellHTML).not.toContain('E2E User');

        const releaseResponse = await request.post(`${controlURL}/release`);
        expect(releaseResponse.ok()).toBe(true);

        const streamedHTML = await readUntil(
            reader,
            decoder,
            (html) => html.includes('E2E User'),
        );

        // backend 解放後、後続チャンクで動的UserMenuが届く
        expect(streamedHTML).toContain('E2E User');
    } finally {
        // assertion失敗時も保留を解除し、後続テストへ未完了レスポンスを残さない
        await request.post(`${controlURL}/release`);
        controller.abort();
        await reader?.cancel().catch(() => undefined);
    }
});

test('[PPR-002] backend障害時もシェルは表示されUserMenuは縮退する', async ({
    authenticatedSession,
    page,
    request,
}) => {
    const outageResponse = await request.post(
        `http://localhost:3101/__e2e/sessions/${authenticatedSession.token}/outage`,
    );
    expect(outageResponse.ok()).toBe(true);

    await page.goto('/dashboard');

    // ページ全体が既定のエラー画面に落ちず、シェルは表示され続ける
    await expect(page.getByRole('heading', { name: 'ホーム' })).toBeVisible();
    await expect(page.getByLabel('ユーザー情報を取得できませんでした')).toBeVisible();
});

test('[PPR-003] ページ内の取得エラーはerror boundaryが受け再試行で復帰する', async ({
    authenticatedSession,
    page,
    request,
}) => {
    const outageResponse = await request.post(
        `http://localhost:3101/__e2e/sessions/${authenticatedSession.token}/outage`,
    );
    expect(outageResponse.ok()).toBe(true);

    await page.goto('/account');
    await expect(
        page.getByText('ページの表示に失敗しました。時間をおいて再度お試しください。'),
    ).toBeVisible();

    const recoverResponse = await request.post(
        `http://localhost:3101/__e2e/sessions/${authenticatedSession.token}/recover`,
    );
    expect(recoverResponse.ok()).toBe(true);

    await page.getByRole('button', { name: '再試行' }).click();

    await expect(page.getByRole('heading', { name: 'アカウント設定' })).toBeVisible();
    await expect(page.getByText('E2E User')).toBeVisible();
});
