// Enforce the dashboard bundle-size budget defined in the Web UI SPA initiative.
// Run this after `npm run build` so the assets exist in internal/ui/dashboard-dist/.
import fs from 'node:fs';
import path from 'node:path';

const DIST_DIR = path.resolve(process.cwd(), '..', '..', 'internal', 'ui', 'dashboard-dist');
const JS_GZ_BUDGET_BYTES = 100 * 1024; // 100 KB
const TOTAL_DIST_BUDGET_BYTES = 250 * 1024; // 250 KB

function dirSize(dir) {
  let total = 0;
  for (const entry of fs.readdirSync(dir, { withFileTypes: true, recursive: true })) {
    if (entry.isFile()) {
      total += fs.statSync(path.join(entry.parentPath ?? dir, entry.name)).size;
    }
  }
  return total;
}

function formatKiB(bytes) {
  return `${(bytes / 1024).toFixed(2)} KiB`;
}

let failed = false;

const jsGzPath = path.join(DIST_DIR, 'assets', 'index.js.gz');
if (!fs.existsSync(jsGzPath)) {
  console.error(`❌ index.js.gz not found at ${jsGzPath}. Run npm run build first.`);
  process.exit(1);
}

const jsGzSize = fs.statSync(jsGzPath).size;
if (jsGzSize > JS_GZ_BUDGET_BYTES) {
  console.error(
    `❌ JS bundle ${formatKiB(jsGzSize)} exceeds ${formatKiB(JS_GZ_BUDGET_BYTES)} budget`,
  );
  failed = true;
} else {
  console.log(`✅ JS bundle ${formatKiB(jsGzSize)} ≤ ${formatKiB(JS_GZ_BUDGET_BYTES)}`);
}

const totalSize = dirSize(DIST_DIR);
if (totalSize > TOTAL_DIST_BUDGET_BYTES) {
  console.error(
    `❌ Total dist ${formatKiB(totalSize)} exceeds ${formatKiB(TOTAL_DIST_BUDGET_BYTES)} budget`,
  );
  failed = true;
} else {
  console.log(`✅ Total dist ${formatKiB(totalSize)} ≤ ${formatKiB(TOTAL_DIST_BUDGET_BYTES)}`);
}

if (failed) {
  process.exit(1);
}
