import { test, expect } from '@playwright/test';
import { spawn, type ChildProcess } from 'node:child_process';
import { mkdtemp, readFile, rm, writeFile } from 'node:fs/promises';
import { createServer } from 'node:net';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import AxeBuilder from '@axe-core/playwright';

const REPO_ROOT = resolve(fileURLToPath(import.meta.url), '..', '..', '..', '..');
const MUARA_URL = process.env.MUARA_URL;

let server: ChildProcess | null = null;
let workdir: string | null = null;
let baseURL: string;

function findFreePort(): Promise<number> {
  return new Promise((resolve, reject) => {
    const srv = createServer();
    srv.listen(0, '127.0.0.1', () => {
      const port = (srv.address() as { port: number }).port;
      srv.close(() => resolve(port));
    });
    srv.on('error', reject);
  });
}

function run(cmd: string, args: string[], cwd: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const child = spawn(cmd, args, { cwd, stdio: 'pipe' });
    let stdout = '';
    let stderr = '';
    child.stdout?.on('data', (d) => { stdout += d; });
    child.stderr?.on('data', (d) => { stderr += d; });
    child.on('close', (code) => {
      if (code === 0) {
        resolve();
      } else {
        reject(new Error(`Command failed with code ${code}: ${cmd} ${args.join(' ')}\nstdout: ${stdout}\nstderr: ${stderr}`));
      }
    });
  });
}

function waitForServer(port: number, timeoutMs = 30000): Promise<void> {
  const deadline = Date.now() + timeoutMs;
  return new Promise((resolve, reject) => {
    const tryConnect = () => {
      const req = spawn('curl', ['-fsS', `http://127.0.0.1:${port}/healthz`]);
      req.on('close', (code) => {
        if (code === 0) {
          resolve();
        } else if (Date.now() > deadline) {
          reject(new Error(`Server did not become healthy within ${timeoutMs}ms`));
        } else {
          setTimeout(tryConnect, 250);
        }
      });
    };
    tryConnect();
  });
}

test.beforeAll(async () => {
  if (MUARA_URL) {
    baseURL = MUARA_URL;
    return;
  }

  const port = await findFreePort();
  workdir = await mkdtemp(join(tmpdir(), 'muara-a11y-'));
  const configPath = join(workdir, 'config.yml');
  const binPath = join(workdir, 'muara');

  await run('go', ['build', '-o', binPath, './cmd/muara'], REPO_ROOT);
  await run(binPath, ['--config', configPath, 'init'], REPO_ROOT);

  let config = await readFile(configPath, 'utf8');
  config = config.replace(/port: 9000/, `port: ${port}`);
  config = config.replace(
    /url: ""/,
    `url: "http://127.0.0.1:${port}/_admin/webhook-receiver"`,
  );
  const providers = ['stripe', 'billplz', 'toyyibpay', 'ipay88', 'fawry', 'senangpay'];
  for (const p of providers) {
    config = config.replace(new RegExp(`${p}:\\n    enabled: false`), `${p}:\n    enabled: true`);
  }
  await writeFile(configPath, config);

  server = spawn(binPath, ['--config', configPath, 'start'], {
    cwd: REPO_ROOT,
    stdio: 'pipe',
  });

  await waitForServer(port);
  baseURL = `http://127.0.0.1:${port}`;
});

test.afterAll(async () => {
  if (server) {
    server.kill('SIGTERM');
    await new Promise<void>((resolve) => {
      server?.on('close', () => resolve());
      setTimeout(() => {
        server?.kill('SIGKILL');
        resolve();
      }, 5000);
    });
  }
  if (workdir) {
    await rm(workdir, { recursive: true, force: true });
  }
});

test.describe('dashboard accessibility smoke test', () => {
  test('skip link bypasses header and tabs', async ({ page }) => {
    await page.goto(`${baseURL}/_admin`);
    await expect(page.locator('main#main-content')).toBeVisible();

    const skipLink = page.locator('.skip-link');
    await expect(skipLink).toHaveCSS('top', '-40px');

    await page.keyboard.press('Tab');
    await expect(skipLink).toBeFocused();
    await expect(skipLink).toHaveCSS('top', '8px');

    await skipLink.click();
    await expect(page.locator('main#main-content')).toBeFocused();
  });

  test('sidebar navigation is keyboard operable', async ({ page }) => {
    await page.goto(`${baseURL}/_admin`);
    const items = page.locator('[role="menuitem"]');
    await expect(items).toHaveCount(3);

    await expect(page.locator('[data-testid="nav-ledger"]')).toHaveAttribute('aria-current', 'page');

    await page.keyboard.press('2');
    await expect(page.locator('[data-testid="nav-webhooks"]')).toHaveAttribute('aria-current', 'page');

    await page.keyboard.press('3');
    await expect(page.locator('[data-testid="nav-settings"]')).toHaveAttribute('aria-current', 'page');

    await page.keyboard.press('1');
    await expect(page.locator('[data-testid="nav-ledger"]')).toHaveAttribute('aria-current', 'page');
  });

  test('theme toggle switches light and dark mode', async ({ page }) => {
    await page.goto(`${baseURL}/_admin`);
    const toggle = page.getByRole('button', { name: /^switch to (light|dark) mode$/i });
    await expect(toggle).toBeVisible();

    await toggle.click();
    await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark');
    await page.screenshot({ path: 'test-results/dashboard-dark.png' });

    await toggle.click();
    await expect(page.locator('html')).toHaveAttribute('data-theme', 'light');
    await page.screenshot({ path: 'test-results/dashboard-light.png' });
  });

  test('dashboard has no critical a11y violations', async ({ page }) => {
    await page.goto(`${baseURL}/_admin`);
    const results = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa'])
      .analyze();
    const critical = results.violations.filter((v) => v.impact === 'critical');
    expect(critical, JSON.stringify(critical, null, 2)).toHaveLength(0);
  });
});
