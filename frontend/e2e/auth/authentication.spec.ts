import { expect, test } from '@playwright/test';

test.describe('未認証時の導線', () => {
    test('保護されたページはログインページへ遷移する', async ({ page }) => {
        await page.goto('/dashboard');

        await expect(page).toHaveURL('/login');
        await expect(page.getByRole('heading', { name: 'Adjusta' })).toBeVisible();
    });
});
