import { test, expect } from '@playwright/test';

const MAILPIT_URL = 'http://127.0.0.1:9035';

test.describe('checkout-store e2e', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('Fawry payment flow completes and redirects to success', async ({ page }) => {
    const name = 'Ahmad Ali';
    const email = 'ahmad.ali@example.com';

    // Landing page → checkout.
    await page.locator('.card .btn-primary').click();
    await expect(page).toHaveURL('/checkout');

    // Fill checkout form.
    await page.locator('input[type="text"]').fill(name);
    await page.locator('input[type="email"]').fill(email);
    await expect(page.locator('button:has-text("Fawry")')).toBeVisible();

    // Submit payment.
    await page.locator('button[type="submit"]').click();

    // OpenMuara Fawry escape page.
    await expect(page).toHaveURL(/\/\_admin\/fawry-escape/);
    await expect(page.locator('text=Simulate Fawry Payment')).toBeVisible();
    await expect(page.locator('text=49.99 EGP')).toBeVisible();

    await page.locator('text=Simulate Paid').click();

    // Back to store success page with paid status.
    await expect(page).toHaveURL(/\/success\?ref=/);
    await expect(page.locator('text=Payment successful!')).toBeVisible();
    await expect(page.locator(`text=${name}`)).toBeVisible();
    await expect(page.locator(`text=${email}`)).toBeVisible();
  });

  test('Stripe payment flow creates session and redirects through pay page', async ({ page }) => {
    const name = 'Siti Aminah';
    const email = 'siti@example.com';

    await page.locator('.card .btn-primary').click();
    await expect(page).toHaveURL('/checkout');

    await page.locator('input[type="text"]').fill(name);
    await page.locator('input[type="email"]').fill(email);

    // Select Stripe.
    await page.locator('.grid button:has-text("Stripe")').click();

    await page.locator('button[type="submit"]').click();

    // OpenMuara Stripe pay page.
    await expect(page).toHaveURL(/\/v1\/checkout\/sessions\/.*\/pay/);
    await expect(page.locator('text=Checkout')).toBeVisible();
    await expect(page.locator('text=49.99 EGP')).toBeVisible();

    await page.locator('text=Confirm payment').click();

    // Redirected back to store.
    await expect(page).toHaveURL(/\/success\?ref=/);
  });

  test('webhook receiver sends confirmation email via Mailpit', async ({ request }) => {
    // Create a checkout through the API so we have a known ref.
    const checkout = await request.post('/api/checkout', {
      data: {
        method: 'fawry',
        name: 'Fatimah Zahra',
        email: 'fatimah@example.com',
      },
    });
    expect(checkout.ok()).toBeTruthy();
    const { ref } = await checkout.json();
    expect(ref).toBeTruthy();

    // Clear Mailpit inbox.
    await request.delete(`${MAILPIT_URL}/api/v1/messages`);

    // Send a Fawry V2-style paid webhook.
    const webhook = await request.post('/webhook', {
      data: {
        merchantRefNumber: ref,
        orderStatus: 'PAID',
      },
    });
    expect(webhook.ok()).toBeTruthy();

    // Wait for the async email delivery and verify Mailpit received it.
    let confirmation: any;
    for (let i = 0; i < 10; i++) {
      const messagesRes = await request.get(`${MAILPIT_URL}/api/v1/messages`);
      expect(messagesRes.ok()).toBeTruthy();
      const mailbox = await messagesRes.json();
      confirmation = mailbox.messages.find((m: any) =>
        m.Subject.includes('Payment confirmed')
      );
      if (confirmation) break;
      await new Promise((r) => setTimeout(r, 200));
    }
    expect(confirmation).toBeTruthy();
    expect(confirmation.To[0].Address).toBe('fatimah@example.com');
  });
});
