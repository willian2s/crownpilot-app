import { spawnSync } from 'node:child_process';
import { executable, repositoryRoot, wait } from './process.mjs';

function docker(args, options = {}) {
  return spawnSync(executable('docker'), args, {
    cwd: repositoryRoot,
    encoding: 'utf8',
    ...options,
  });
}

const version = docker(['version', '--format', '{{.Server.Version}}']);
if (version.error || version.status !== 0) {
  console.error('Container smoke blocked: Docker daemon/CLI is unavailable.');
  process.exit(2);
}

const image = process.env.CROWNPILOT_CONTAINER_IMAGE ?? 'crownpilot-api:bootstrap';
const container = `crownpilot-api-smoke-${process.pid}`;
const port = process.env.CROWNPILOT_SMOKE_PORT ?? '18080';

try {
  const build = docker(['build', '--tag', image, '.'], { stdio: 'inherit' });
  if (build.status !== 0) {
    process.exit(build.status ?? 1);
  }

  const start = docker(['run', '--detach', '--name', container, '--publish', `${port}:8080`, image], {
    stdio: 'inherit',
  });
  if (start.status !== 0) {
    process.exit(start.status ?? 1);
  }

  let healthy = false;
  for (let attempt = 0; attempt < 40; attempt += 1) {
    try {
      const response = await fetch(`http://127.0.0.1:${port}/health/live`);
      if (response.ok && (await response.text()) === 'Healthy') {
        healthy = true;
        break;
      }
    } catch {
      // Container may still be starting.
    }
    await wait(250);
  }

  if (!healthy) {
    throw new Error('Container did not pass /health/live before timeout.');
  }

  console.log('Container smoke passed: /health/live');
} finally {
  docker(['rm', '--force', container], { stdio: 'inherit' });
}
