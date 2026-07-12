// Automated WCAG color-contrast regression check for the dashboard.
// Run after the Go server (and built dashboard assets) are available.
// Usage: node scripts/a11y-contrast-check.js [url]
// Defaults to MUARA_URL env var or http://127.0.0.1:9000
// If no URL is reachable, the script builds and starts a temporary OpenMuara server.

import { chromium } from 'playwright';
import AxeBuilder from '@axe-core/playwright';
import { spawn } from 'node:child_process';
import { mkdtemp, readFile, rm, writeFile } from 'node:fs/promises';
import { createServer } from 'node:net';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const REPO_ROOT = resolve(fileURLToPath(import.meta.url), '..', '..', '..', '..');

function findFreePort() {
  return new Promise((resolve, reject) => {
    const srv = createServer();
    srv.listen(0, '127.0.0.1', () => {
      const port = (srv.address()).port;
      srv.close(() => resolve(port));
    });
    srv.on('error', reject);
  });
}

function run(cmd, args, cwd) {
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

function waitForServer(port, timeoutMs = 30000) {
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

async function startLocalServer() {
  const port = await findFreePort();
  const workdir = await mkdtemp(join(tmpdir(), 'muara-contrast-'));
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

  const server = spawn(binPath, ['--config', configPath, 'start'], {
    cwd: REPO_ROOT,
    stdio: 'pipe',
  });

  await waitForServer(port);
  return { url: `http://127.0.0.1:${port}`, server, workdir };
}

async function setTheme(page, theme) {
  await page.evaluate((t) => {
    localStorage.setItem('muara-theme', t);
    document.documentElement.setAttribute('data-theme', t);
    const meta = document.getElementById('theme-color');
    if (meta) meta.content = t === 'dark' ? '#0f172a' : '#f8fafc';
    window.dispatchEvent(new StorageEvent('storage', { key: 'muara-theme', newValue: t }));
  }, theme);
  await page.waitForFunction((t) => document.documentElement.getAttribute('data-theme') === t, theme);
  await page.waitForTimeout(150);
}

async function checkUrl(context, url) {
  const page = await context.newPage();
  await page.goto(url, { waitUntil: 'networkidle' });

  // Force light mode for a deterministic baseline, then check dark mode.
  await setTheme(page, 'light');
  await page.waitForFunction(() => document.documentElement.getAttribute('data-theme') === 'light');

  const lightResults = await new AxeBuilder({ page })
    .withTags(['wcag2aa'])
    .disableRules(['bypass', 'html-has-lang', 'landmark-one-main', 'page-has-heading-one'])
    .analyze();

  const contrastViolations = lightResults.violations.filter((v) => v.id === 'color-contrast');
  if (contrastViolations.length > 0) {
    console.error('❌ WCAG AA contrast violations found (light theme):');
    for (const v of contrastViolations) {
      console.error(`  ${v.help} (${v.nodes.length} node(s))`);
      for (const node of v.nodes.slice(0, 5)) {
        console.error(`    - ${node.target.join(' ')}`);
      }
    }
    return false;
  }
  console.log('✅ No WCAG AA contrast violations (light theme)');

  await setTheme(page, 'dark');

  const darkResults = await new AxeBuilder({ page })
    .withTags(['wcag2aa'])
    .disableRules(['bypass', 'html-has-lang', 'landmark-one-main', 'page-has-heading-one'])
    .analyze();

  const darkContrast = darkResults.violations.filter((v) => v.id === 'color-contrast');
  if (darkContrast.length > 0) {
    console.error('❌ WCAG AA contrast violations found (dark theme):');
    for (const v of darkContrast) {
      console.error(`  ${v.help} (${v.nodes.length} node(s))`);
      for (const node of v.nodes.slice(0, 5)) {
        console.error(`    - ${node.target.join(' ')}`);
      }
    }
    return false;
  }
  console.log('✅ No WCAG AA contrast violations (dark theme)');

  return true;
}

async function main() {
  let url = process.argv[2] || process.env.MUARA_URL || 'http://127.0.0.1:9000';
  url = url.replace(/\/$/, '');

  let server = null;
  let workdir = null;

  try {
    const probe = await fetch(`${url}/healthz`, { signal: AbortSignal.timeout(2000) }).catch(() => null);
    if (!probe || !probe.ok) {
      console.log('No server reachable; starting a temporary OpenMuara server...');
      const local = await startLocalServer();
      url = local.url;
      server = local.server;
      workdir = local.workdir;
    }

    const browser = await chromium.launch();
    const context = await browser.newContext({ viewport: { width: 1280, height: 720 } });
    const ok = await checkUrl(context, `${url}/_admin`);
    await browser.close();

    if (!ok) {
      process.exit(1);
    }
  } finally {
    if (server) {
      server.kill('SIGTERM');
      await new Promise((resolve) => {
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
  }
}

main();
