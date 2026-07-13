import { expect, test } from '@playwright/test';

test('[PUBLIC-002] ログインページを表示できる', async ({ page }) => {
    await page.goto('/login');

    await expect(page.getByRole('heading', { name: 'Adjusta' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Googleでログイン' })).toBeVisible();
});
