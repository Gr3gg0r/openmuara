const fs = require('fs');
const path = require('path');
const markdownLinkCheck = require('markdown-link-check');

const rootDir = path.resolve(__dirname, '..', '..');
const configPath = path.join(rootDir, '.markdown-link-check.json');

function loadConfig() {
  const raw = fs.readFileSync(configPath, 'utf8');
  const config = JSON.parse(raw);

  // The markdown-link-check API expects timeouts in milliseconds.
  if (typeof config.timeout === 'string') {
    config.timeout = ms(config.timeout);
  }
  if (typeof config.fallbackRetryDelay === 'string') {
    config.fallbackRetryDelay = ms(config.fallbackRetryDelay);
  }

  return config;
}

function ms(value) {
  const match = /^([0-9]+)\s*(ms|s|m|h)?$/i.exec(value);
  if (!match) return 0;

  const amount = parseInt(match[1], 10);
  const unit = (match[2] || 'ms').toLowerCase();
  switch (unit) {
    case 'ms':
      return amount;
    case 's':
      return amount * 1000;
    case 'm':
      return amount * 60 * 1000;
    case 'h':
      return amount * 60 * 60 * 1000;
    default:
      return amount;
  }
}

function collectFiles() {
  const files = new Set();

  function addFile(relative) {
    const resolved = path.resolve(rootDir, relative);
    if (!fs.existsSync(resolved) || !fs.statSync(resolved).isFile()) {
      console.error(`File not found: ${relative}`);
      return false;
    }
    files.add(resolved);
    return true;
  }

  function addDirectory(relative) {
    const resolved = path.resolve(rootDir, relative);
    if (!fs.existsSync(resolved) || !fs.statSync(resolved).isDirectory()) {
      return;
    }

    const entries = fs.readdirSync(resolved, { recursive: true });
    for (const entry of entries) {
      const fullPath = path.join(resolved, entry);
      if (fs.statSync(fullPath).isFile() && fullPath.endsWith('.md')) {
        files.add(fullPath);
      }
    }
  }

  addFile('README.md');
  addFile('AGENTS.md');
  addDirectory('docs');
  addDirectory('runbooks');
  addDirectory('.github');

  return Array.from(files).sort();
}

function checkFile(file, config) {
  return new Promise((resolve) => {
    const markdown = fs.readFileSync(file, 'utf8');
    const fileConfig = { ...config, baseUrl: `file://${path.dirname(file)}/` };
    markdownLinkCheck(markdown, fileConfig, (err, results) => {
      if (err) {
        console.error(`\n${path.relative(rootDir, file)}:`);
        console.error(`  error: ${err.message}`);
        resolve(1);
        return;
      }

      const dead = results.filter((r) => r.status === 'dead');
      if (dead.length > 0) {
        console.error(`\n${path.relative(rootDir, file)}:`);
        for (const result of dead) {
          const detail = result.statusCode
            ? `status ${result.statusCode}`
            : result.err?.message || 'unknown error';
          console.error(`  dead: ${result.link} (${detail})`);
        }
        resolve(dead.length);
        return;
      }

      resolve(0);
    });
  });
}

async function main() {
  const config = loadConfig();
  const files = collectFiles();

  if (files.length === 0) {
    console.log('No markdown files to check.');
    process.exit(0);
  }

  let failures = 0;
  for (const file of files) {
    failures += await checkFile(file, config);
  }

  if (failures > 0) {
    console.error(`\n${failures} dead link(s) found.`);
    process.exit(1);
  }

  console.log(`Checked ${files.length} file(s). All links passed.`);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
