import { test, expect } from '@playwright/test';

test('homepage returns 200', async ({ page }) => {
	const response = await page.goto('/');
	expect(response?.status()).toBe(200);
});

test('health endpoint returns ok', async ({ request }) => {
	const response = await request.get('/api/health');
	expect(response.status()).toBe(200);

	const body = await response.json();
	expect(body.status).toBe('ok');
});
