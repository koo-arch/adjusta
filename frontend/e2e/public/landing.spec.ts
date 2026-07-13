import { expect, test } from '@playwright/test';

test('トップページを表示できる', async ({ page }) => {
    await page.goto('/');

    await expect(
        page.getByRole('heading', { name: '日程調整をもっとシンプルに' }),
    ).toBeVisible();
    await page.getByRole('button', { name: '今すぐ始める' }).click();

    await expect(page).toHaveURL('/login');
});
