import { expect, test } from '../fixtures/auth';
import type { APIRequestContext } from '@playwright/test';

test.describe.configure({ mode: 'serial' });

const configureCandidateSync = async (
    request: APIRequestContext,
    mode: 'off' | 'on' | 'fail-update',
) => {
    const response = await request.post(`http://localhost:3101/__e2e/candidate-sync/${mode}`);
    expect(response.ok()).toBe(true);
};

const configureCalendarSettings = async (
    request: APIRequestContext,
    mode: 'reset' | 'fail-update' = 'reset',
) => {
    const response = await request.post(`http://localhost:3101/__e2e/calendar-settings/${mode}`);
    expect(response.ok()).toBe(true);
};

test('[ACCOUNT-001] 候補日程の同期は初期状態でOFFになる', async ({
    authenticatedSession: _authenticatedSession,
    page,
    request,
}) => {
    await configureCandidateSync(request, 'off');
    await page.goto('/account');

    await expect(page.getByRole('heading', { name: 'アカウント設定' })).toBeVisible();
    await expect(page.getByText('カレンダー設定', { exact: true })).toBeVisible();
    await expect(page.getByText('専用カレンダーに仮予定を表示します。')).toBeVisible();
    await expect(
        page.getByRole('switch', { name: '候補日程の Google Calendar 同期' }),
    ).not.toBeChecked();
});

test('[ACCOUNT-002] 候補日程の同期をONにして専用カレンダーを作成できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
    request,
}) => {
    await configureCandidateSync(request, 'off');
    await page.goto('/account');
    const candidateSync = page.getByRole('switch', {
        name: '候補日程の Google Calendar 同期',
    });
    await expect(candidateSync).not.toBeChecked();

    const updateRequestPromise = page.waitForRequest(
        (apiRequest) =>
            apiRequest.method() === 'PUT' &&
            apiRequest.url().endsWith('/api/calendar-settings/candidate-sync'),
    );
    await candidateSync.click();

    const updateRequest = await updateRequestPromise;
    expect(updateRequest.postDataJSON()).toEqual({ enabled: true });
    await expect(candidateSync).toBeChecked();
    await expect(page.getByText('Adjusta 調整用', { exact: true })).toBeVisible();
});

test('[ACCOUNT-003] 候補日程の同期をOFFにしても専用カレンダーを維持する', async ({
    authenticatedSession: _authenticatedSession,
    page,
    request,
}) => {
    await configureCandidateSync(request, 'on');
    await page.goto('/account');
    const candidateSync = page.getByRole('switch', {
        name: '候補日程の Google Calendar 同期',
    });
    await expect(candidateSync).toBeChecked();

    const updateRequestPromise = page.waitForRequest(
        (apiRequest) =>
            apiRequest.method() === 'PUT' &&
            apiRequest.url().endsWith('/api/calendar-settings/candidate-sync'),
    );
    await candidateSync.click();

    const updateRequest = await updateRequestPromise;
    expect(updateRequest.postDataJSON()).toEqual({ enabled: false });
    await expect(candidateSync).not.toBeChecked();
    await expect(page.getByText('Adjusta 調整用', { exact: true })).toBeVisible();
});

test('[ACCOUNT-004] 同期設定の更新に失敗した場合はOFFへ戻す', async ({
    authenticatedSession: _authenticatedSession,
    page,
    request,
}) => {
    await configureCandidateSync(request, 'fail-update');
    await page.goto('/account');
    const candidateSync = page.getByRole('switch', {
        name: '候補日程の Google Calendar 同期',
    });
    await expect(candidateSync).not.toBeChecked();

    await candidateSync.click();

    await expect(page.getByText('候補日程の同期設定を更新できませんでした')).toBeVisible();
    await expect(candidateSync).not.toBeChecked();
});

test('[ACCOUNT-005] 確定予定の登録先から候補用と祝日カレンダーを除外する', async ({
    authenticatedSession: _authenticatedSession,
    page,
    request,
}) => {
    await configureCandidateSync(request, 'on');
    await configureCalendarSettings(request);
    await page.goto('/account');

    await expect(page.getByRole('radio', { name: 'メインカレンダー' })).toBeChecked();
    await expect(page.getByRole('radio', { name: 'チームカレンダー' })).toBeVisible();
    await expect(page.getByRole('radio', { name: 'Adjusta 調整用' })).toHaveCount(0);
    await expect(page.getByRole('radio', { name: '日本の祝日' })).toHaveCount(0);
});

test('[ACCOUNT-006] 確定予定の登録先を変更できる', async ({
    authenticatedSession: _authenticatedSession,
    page,
    request,
}) => {
    await configureCandidateSync(request, 'off');
    await configureCalendarSettings(request);
    await page.goto('/account');
    const teamCalendar = page.getByRole('radio', { name: 'チームカレンダー' });

    const updateRequestPromise = page.waitForRequest(
        (apiRequest) =>
            apiRequest.method() === 'PATCH' &&
            apiRequest.url().endsWith('/api/user-calendars/calendar-reference'),
    );
    await teamCalendar.click();

    const updateRequest = await updateRequestPromise;
    expect(updateRequest.postDataJSON()).toEqual({ role: 'primary' });
    await expect(teamCalendar).toBeChecked();
    await expect(page.getByRole('radio', { name: 'メインカレンダー' })).not.toBeChecked();
});

test('[ACCOUNT-007] 予定確認に表示するカレンダーをOFFにできる', async ({
    authenticatedSession: _authenticatedSession,
    page,
    request,
}) => {
    await configureCandidateSync(request, 'off');
    await configureCalendarSettings(request);
    await page.goto('/account');
    const primaryRow = page.getByRole('listitem').filter({ hasText: 'メインカレンダー' });
    const visibleSwitch = primaryRow.getByRole('switch');
    await expect(visibleSwitch).toBeChecked();

    const updateRequestPromise = page.waitForRequest(
        (apiRequest) =>
            apiRequest.method() === 'PATCH' &&
            apiRequest.url().endsWith('/api/user-calendars/calendar-primary'),
    );
    await visibleSwitch.click();

    const updateRequest = await updateRequestPromise;
    expect(updateRequest.postDataJSON()).toEqual({ is_visible: false });
    await expect(visibleSwitch).not.toBeChecked();
});

test('[ACCOUNT-008] カレンダー表示設定の更新に失敗した場合はONへ戻す', async ({
    authenticatedSession: _authenticatedSession,
    page,
    request,
}) => {
    await configureCandidateSync(request, 'off');
    await configureCalendarSettings(request, 'fail-update');
    await page.goto('/account');
    const primaryRow = page.getByRole('listitem').filter({ hasText: 'メインカレンダー' });
    const visibleSwitch = primaryRow.getByRole('switch');
    await expect(visibleSwitch).toBeChecked();

    await visibleSwitch.click();

    await expect(page.getByText('カレンダー設定の更新に失敗しました')).toBeVisible();
    await expect(visibleSwitch).toBeChecked();
});
