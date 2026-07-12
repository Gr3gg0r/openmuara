import { test, expect } from '@playwright/test';

test.describe('OpenMuara dashboard', () => {
  test('loads in hardened mode with credentials in the URL', async ({ page }) => {
    const errors: string[] = [];
    page.on('pageerror', (err) => errors.push(err.message));
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        errors.push(msg.text());
      }
    });

    await page.goto('http://admin:testpass123@127.0.0.1:59005/_admin/');
    await expect(page.locator('text=OpenMuara Dashboard')).toBeVisible();
    await expect(page.locator('text=Getting started')).toBeVisible();

    // The Fetch API rejects URLs that contain credentials. The SPA must strip
    // them before calling fetch, otherwise every API request throws and the
    // dashboard shows an error banner.
    const fetchCredentialError = errors.find((e) =>
      e.includes('Request cannot be constructed from a URL that includes credentials'),
    );
    expect(fetchCredentialError).toBeUndefined();
  });
});
