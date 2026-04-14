import { test, expect } from '@playwright/test';

test.describe('theme switching', () => {
	test('light theme is default', async ({ page }) => {
		await page.goto('/auth');
		const htmlClass = await page.locator('html').getAttribute('class');
		expect(htmlClass ?? '').not.toContain('dark');
	});

	test('dark theme persists after reload', async ({ page }) => {
		await page.goto('/auth');
		await page.evaluate(() => {
			localStorage.setItem('pulse_theme', 'dark');
		});
		await page.reload();
		const htmlClass = await page.locator('html').getAttribute('class');
		expect(htmlClass).toContain('dark');
	});

	test('light theme persists after reload', async ({ page }) => {
		await page.goto('/auth');
		await page.evaluate(() => {
			localStorage.setItem('pulse_theme', 'light');
		});
		await page.reload();
		const htmlClass = await page.locator('html').getAttribute('class');
		expect(htmlClass ?? '').not.toContain('dark');
	});
});
