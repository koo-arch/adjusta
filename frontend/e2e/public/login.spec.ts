import { expect, test } from '@playwright/test';

test('[PUBLIC-002] ログインページを表示できる', async ({ page }) => {
    await page.goto('/login');

    await expect(page.getByRole('heading', { name: 'Adjusta' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Googleでログイン' })).toBeVisible();
});

test('[PUBLIC-003] Googleログインをsame-origin proxy経由で完了できる', async ({ page }) => {
    await page.goto('/login');
    await page.getByRole('button', { name: 'Googleでログイン' }).click();

    await expect(page).toHaveURL('/dashboard');
    await expect(page.getByRole('heading', { name: 'ホーム' })).toBeVisible();
});
