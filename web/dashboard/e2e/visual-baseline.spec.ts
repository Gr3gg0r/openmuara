import { test, expect, type Page } from '@playwright/test';
import { spawn, type ChildProcess } from 'node:child_process';
import { mkdtemp, readFile, rm, writeFile } from 'node:fs/promises';
import { createServer } from 'node:net';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const REPO_ROOT = resolve(fileURLToPath(import.meta.url), '..', '..', '..', '..');
const MUARA_URL = process.env.MUARA_URL;
const BASELINE_DIR = resolve(REPO_ROOT, 'docs', 'initiatives', 'openmuara-bug-hunt', 'findings', 'visual-baseline');

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
        reject(new Error(`Command failed with code ${code}: ${cmd} ${args.join(' ')}
stdout: ${stdout}
stderr: ${stderr}`));
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
  workdir = await mkdtemp(join(tmpdir(), 'muara-visual-'));
  const configPath = join(workdir, 'config.yml');
  const dbPath = join(workdir, 'ledger.db');
  const binPath = join(workdir, 'muara');

  await run('go', ['build', '-o', binPath, './cmd/muara'], REPO_ROOT);
  await run(binPath, ['--config', configPath, 'init'], REPO_ROOT);

  let config = await readFile(configPath, 'utf8');
  config = config.replace(/port: 9000/, `port: ${port}`);
  config = config.replace(/path: \.muara\/data\/ledger\.db/, `path: ${dbPath}`);
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

async function setTheme(page: Page, theme: 'light' | 'dark'): Promise<void> {
  await page.goto(`${baseURL}/_admin`, { waitUntil: 'networkidle' });
  await page.evaluate((t) => {
    localStorage.setItem('muara-theme', t);
  }, theme);
  await page.reload({ waitUntil: 'networkidle' });
}

test.describe('visual baseline', () => {
  for (const theme of ['light', 'dark'] as const) {
    test(`dashboard views match baseline (${theme})`, async ({ page }) => {
      await setTheme(page, theme);

      await page.goto(`${baseURL}/_admin`, { waitUntil: 'networkidle' });
      await expect(page).toHaveScreenshot(`ledger-default-${theme}.png`, { fullPage: true });

      await page.goto(`${baseURL}/_admin?view=webhooks`, { waitUntil: 'networkidle' });
      await expect(page).toHaveScreenshot(`webhooks-view-${theme}.png`, { fullPage: true });

      await page.goto(`${baseURL}/_admin?view=settings`, { waitUntil: 'networkidle' });
      await expect(page).toHaveScreenshot(`settings-view-${theme}.png`, { fullPage: true });

      await page.goto(`${baseURL}/_admin?view=settings&provider=fawry`, { waitUntil: 'networkidle' });
      await expect(page).toHaveScreenshot(`provider-detail-fawry-${theme}.png`, { fullPage: true });

      await page.goto(`${baseURL}/_admin?view=settings&provider=stripe`, { waitUntil: 'networkidle' });
      await expect(page).toHaveScreenshot(`provider-detail-stripe-${theme}.png`, { fullPage: true });
    });
  }
});
