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

  // Sem CROWNPILOT_ENVIRONMENT a imagem precisa falhar fechado, antes do listener.
  const unconfigured = docker(['run', '--rm', image]);
  if (unconfigured.status === 0 || !unconfigured.stderr.includes('CROWNPILOT_ENVIRONMENT')) {
    throw new Error('Container without CROWNPILOT_ENVIRONMENT did not fail closed.');
  }

  // Production: mesma imagem que será promovida; OpenAPI e /docs desligados.
  const start = docker(
    [
      'run',
      '--detach',
      '--name',
      container,
      '--env',
      'CROWNPILOT_ENVIRONMENT=Production',
      '--publish',
      `${port}:8080`,
      image,
    ],
    { stdio: 'inherit' },
  );
  if (start.status !== 0) {
    process.exit(start.status ?? 1);
  }

  const baseUrl = `http://127.0.0.1:${port}`;
  let healthy = false;
  for (let attempt = 0; attempt < 40; attempt += 1) {
    try {
      const response = await fetch(`${baseUrl}/health/live`);
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

  const openapi = await fetch(`${baseUrl}/openapi/v1.json`);
  if (openapi.status !== 404) {
    throw new Error(`Production container exposed /openapi/v1.json (HTTP ${openapi.status}).`);
  }

  console.log('Container smoke passed: fail-closed start, /health/live, OpenAPI hidden in Production');
} finally {
  docker(['rm', '--force', container], { stdio: 'inherit' });
}
