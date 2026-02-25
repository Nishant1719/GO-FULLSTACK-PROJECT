#!/usr/bin/env node
/**
 * Early health checks - verifies all services are reachable.
 * Run before development or in CI.
 */

const services = [
  { name: 'Go Domain API', url: process.env.GO_API_URL || 'http://localhost:8080', path: '/ping' },
  { name: 'BFF Node', url: process.env.BFF_URL || 'http://localhost:3000', path: '/health' },
];

async function check(url, path) {
  try {
    const res = await fetch(`${url}${path}`);
    const ok = res.ok;
    const data = res.ok ? await res.json() : {};
    return { ok, status: res.status, data };
  } catch (e) {
    return { ok: false, error: e.message };
  }
}

async function main() {
  console.log('Running early health checks...\n');
  let allOk = true;

  for (const svc of services) {
    const result = await check(svc.url, svc.path);
    const status = result.ok ? '✓ OK' : '✗ FAIL';
    console.log(`${svc.name} (${svc.url}${svc.path}): ${status}`);
    if (!result.ok) {
      allOk = false;
      if (result.error) console.log(`  Error: ${result.error}`);
      else console.log(`  Status: ${result.status}`);
    }
  }

  console.log('');
  process.exit(allOk ? 0 : 1);
}

main();
