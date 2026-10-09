import { spawn, spawnSync } from 'node:child_process';
import { once } from 'node:events';
import { mkdtemp, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

export const frontendDirectory = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
export const repositoryRoot = path.resolve(frontendDirectory, '..');
export const apiPackage = './cmd/crownpilot-api';

export function executable(command) {
  return process.platform === 'win32' ? `${command}.cmd` : command;
}

export function run(command, args, options = {}) {
  const result = spawnSync(executable(command), args, {
    cwd: repositoryRoot,
    stdio: 'inherit',
    ...options,
  });

  if (result.error) {
    throw result.error;
  }

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

export async function wait(milliseconds) {
  await new Promise((resolve) => setTimeout(resolve, milliseconds));
}

export async function waitForHttp(url, attempts = 40) {
  for (let attempt = 0; attempt < attempts; attempt += 1) {
    try {
      const response = await fetch(url);
      if (response.ok) {
        return response;
      }
    } catch {
      // The process may still be starting; retry within the bounded window.
    }

    await wait(250);
  }

  throw new Error(`Timed out waiting for ${url}`);
}

// Compila a API Go num diretório temporário e executa o binário diretamente.
// `go run` é evitado de propósito: ele inicia um processo filho que um SIGTERM
// enviado ao próprio `go run` não encerra de forma confiável.
export async function withLocalApi(check) {
  const port = process.env.CROWNPILOT_SMOKE_PORT ?? '5089';
  const baseUrl = `http://127.0.0.1:${port}`;
  const buildDirectory = await mkdtemp(path.join(tmpdir(), 'crownpilot-api-'));
  const binary = path.join(buildDirectory, process.platform === 'win32' ? 'crownpilot-api.exe' : 'crownpilot-api');

  try {
    run('go', ['build', '-o', binary, apiPackage]);

    const child = spawn(binary, [], {
      cwd: repositoryRoot,
      env: {
        ...process.env,
        CROWNPILOT_ENVIRONMENT: 'Local',
        PORT: port,
      },
      // stdout (logs JSON) não é lido; um pipe sem leitura poderia encher o buffer.
      stdio: ['ignore', 'ignore', 'pipe'],
    });
    let stderr = '';
    child.stderr?.on('data', (chunk) => {
      stderr += chunk.toString();
    });

    try {
      await waitForHttp(`${baseUrl}/health/live`);
      await check(baseUrl);
    } catch (error) {
      const details = stderr.trim();
      throw new Error(
        `${error instanceof Error ? error.message : String(error)}${details ? `\n${details}` : ''}`,
        { cause: error },
      );
    } finally {
      if (child.exitCode === null) {
        child.kill('SIGTERM');
        await Promise.race([once(child, 'exit'), wait(5_000)]);
        if (child.exitCode === null) {
          child.kill('SIGKILL');
        }
      }
    }
  } finally {
    await rm(buildDirectory, { recursive: true, force: true });
  }
}
