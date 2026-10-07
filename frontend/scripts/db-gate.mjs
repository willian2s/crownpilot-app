import { spawnSync } from 'node:child_process';
import { existsSync, readdirSync, readFileSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const operation = process.argv[2] ?? 'unknown';
const frontendDirectory = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const repositoryRoot = path.resolve(frontendDirectory, '..');
const composeFile = path.join(repositoryRoot, 'database', 'docker-compose.yml');
const securityDirectory = path.join(repositoryRoot, 'database', 'security');
const fixturesDirectory = path.join(repositoryRoot, 'database', 'fixtures');
const infrastructureProject = path.join(
  repositoryRoot,
  'src',
  'Infrastructure',
  'Infrastructure.csproj',
);
const composeProject = 'crownpilot-local';
const databaseName = 'crownpilot_local';
const databaseUser = 'postgres';
const databasePassword = 'postgres';
const defaultDatabasePort = 54322;
const preparedSchemaMigration = '20261007000000_PrepareCrownPilotSchema';

class BlockedDatabaseGateError extends Error {
  exitCode = 2;
}

class CommandError extends Error {
  constructor(message, exitCode = 1) {
    super(message);
    this.exitCode = exitCode;
  }
}

function failClosed(message) {
  throw new BlockedDatabaseGateError(`db:${operation} blocked: ${message}`);
}

function normalizeEnvironment(value) {
  if (!value || value.toLowerCase() === 'development' || value.toLowerCase() === 'local') {
    return 'Local';
  }

  const normalized = value[0].toUpperCase() + value.slice(1).toLowerCase();
  if (['Preview', 'Staging', 'Production'].includes(normalized)) {
    return normalized;
  }

  failClosed(`unsupported runtime environment '${value}'.`);
}

function assertLocalEnvironment() {
  const configuredEnvironments = [
    process.env.CROWNPILOT__ENVIRONMENT,
    process.env.CROWNPILOT_ENVIRONMENT,
    process.env.ASPNETCORE_ENVIRONMENT,
  ].filter(Boolean);

  for (const value of configuredEnvironments) {
    if (normalizeEnvironment(value) !== 'Local') {
      failClosed(
        `database commands are local-only; refusing target environment '${normalizeEnvironment(value)}'.`,
      );
    }
  }
}

function parseConnectionString(connectionString) {
  const values = new Map();
  for (const part of connectionString.split(';')) {
    const separator = part.indexOf('=');
    if (separator < 1) {
      continue;
    }

    const key = part.slice(0, separator).trim().toLowerCase();
    const value = part.slice(separator + 1).trim().replace(/^"|"$/g, '');
    values.set(key, value);
  }

  return values;
}

function localPort() {
  const configuredPort = process.env.CROWNPILOT_DB_PORT ?? String(defaultDatabasePort);
  const port = Number(configuredPort);
  if (!Number.isInteger(port) || port < 1 || port > 65535) {
    failClosed('CROWNPILOT_DB_PORT must be an integer between 1 and 65535.');
  }

  return port;
}

function databaseConnectionString() {
  const port = localPort();
  return (
    process.env.CROWNPILOT__DATABASE__CONNECTIONSTRING ??
    process.env.DATABASE_CONNECTION_STRING ??
    `Host=127.0.0.1;Port=${port};Database=${databaseName};Username=${databaseUser};Password=${databasePassword};SSL Mode=Disable;Application Name=CrownPilot.DbScripts`
  );
}

function assertLocalDatabaseTarget() {
  const configured = parseConnectionString(databaseConnectionString());
  const host = configured.get('host') ?? configured.get('server');
  const port = Number(configured.get('port') ?? defaultDatabasePort);
  const database = configured.get('database') ?? configured.get('initial catalog');
  const username = configured.get('username') ?? configured.get('user id') ?? configured.get('user');

  if (
    !['localhost', '127.0.0.1', '::1'].includes(host?.toLowerCase()) ||
    port !== localPort() ||
    database !== databaseName ||
    username !== databaseUser
  ) {
    failClosed(
      'connection target is not the disposable local PostgreSQL instance; remote targets are never accepted by db:*.',
    );
  }

  return databaseConnectionString();
}

function commandEnvironment(overrides = {}) {
  return {
    ...process.env,
    ASPNETCORE_ENVIRONMENT: 'Development',
    CROWNPILOT__ENVIRONMENT: 'Local',
    CROWNPILOT_DB_PORT: String(localPort()),
    ...overrides,
  };
}

function run(command, args, { env = process.env, input, capture = false } = {}) {
  const result = spawnSync(command, args, {
    cwd: repositoryRoot,
    env,
    encoding: 'utf8',
    input,
    stdio: capture ? ['ignore', 'pipe', 'pipe'] : input !== undefined ? ['pipe', 'inherit', 'inherit'] : 'inherit',
  });

  if (result.error) {
    if (command === 'docker') {
      throw new BlockedDatabaseGateError(
        `db:${operation} blocked: Docker CLI is unavailable. Install Docker Desktop/Engine and retry.`,
      );
    }

    throw new CommandError(`${command} could not be started: ${result.error.message}`);
  }

  if (result.status !== 0) {
    throw new CommandError(
      `${command} ${args.join(' ')} failed with exit code ${result.status ?? 1}.`,
      result.status ?? 1,
    );
  }

  return capture ? (result.stdout ?? '').trim() : '';
}

function compose(args, options = {}) {
  return run(
    'docker',
    ['compose', '--project-name', composeProject, '--file', composeFile, ...args],
    { env: commandEnvironment(), ...options },
  );
}

function assertDockerDaemon() {
  const result = spawnSync('docker', ['info', '--format', '{{.ServerVersion}}'], {
    cwd: repositoryRoot,
    env: commandEnvironment(),
    encoding: 'utf8',
    stdio: ['ignore', 'pipe', 'pipe'],
  });

  if (result.error || result.status !== 0) {
    throw new BlockedDatabaseGateError(
      'db gate requires a running Docker daemon; Docker CLI is installed but the daemon is unavailable.',
    );
  }
}

function sqlFiles(directory) {
  if (!existsSync(directory)) {
    return [];
  }

  return readdirSync(directory)
    .filter((file) => file.endsWith('.sql'))
    .sort()
    .map((file) => path.join(directory, file));
}

function sqlContents(files) {
  return files.map((file) => `-- ${path.relative(repositoryRoot, file)}\n${readFileSync(file, 'utf8')}`).join('\n');
}

function runPsql(sql) {
  return compose(
    [
      'exec',
      '--no-TTY',
      'postgres',
      'psql',
      '--username',
      databaseUser,
      '--dbname',
      databaseName,
      '--set',
      'ON_ERROR_STOP=1',
      '--file=-',
    ],
    { input: sql },
  );
}

function queryPsql(sql) {
  return compose(
    [
      'exec',
      '--no-TTY',
      'postgres',
      'psql',
      '--username',
      databaseUser,
      '--dbname',
      databaseName,
      '--tuples-only',
      '--no-align',
      '--set',
      'ON_ERROR_STOP=1',
      '--command',
      sql,
    ],
    { capture: true },
  );
}

async function waitForDatabase() {
  for (let attempt = 0; attempt < 40; attempt += 1) {
    const result = spawnSync(
      'docker',
      [
        'compose',
        '--project-name',
        composeProject,
        '--file',
        composeFile,
        'exec',
        '--no-TTY',
        'postgres',
        'pg_isready',
        '--username',
        databaseUser,
        '--dbname',
        databaseName,
      ],
      {
        cwd: repositoryRoot,
        env: commandEnvironment(),
        encoding: 'utf8',
        stdio: ['ignore', 'pipe', 'pipe'],
      },
    );

    if (!result.error && result.status === 0) {
      return;
    }

    await new Promise((resolve) => setTimeout(resolve, 250));
  }

  throw new BlockedDatabaseGateError('PostgreSQL container did not become ready within 10 seconds.');
}

function assertDatabaseReachable() {
  assertDockerDaemon();
  const result = spawnSync(
    'docker',
    [
      'compose',
      '--project-name',
      composeProject,
      '--file',
      composeFile,
      'exec',
      '--no-TTY',
      'postgres',
      'pg_isready',
      '--username',
      databaseUser,
      '--dbname',
      databaseName,
    ],
    {
      cwd: repositoryRoot,
      env: commandEnvironment(),
      encoding: 'utf8',
      stdio: ['ignore', 'pipe', 'pipe'],
    },
  );

  if (result.error || result.status !== 0) {
    throw new BlockedDatabaseGateError('local PostgreSQL is not running; run db:start first.');
  }
}

async function start() {
  assertDockerDaemon();
  compose(['up', '--detach', 'postgres']);
  await waitForDatabase();
  console.log(`Local PostgreSQL ready on 127.0.0.1:${localPort()}.`);
}

function stop() {
  assertDockerDaemon();
  compose(['down', '--remove-orphans']);
  console.log('Local PostgreSQL stopped.');
}

function migrate() {
  assertLocalDatabaseTarget();
  assertDatabaseReachable();
  run(
    'dotnet',
    [
      'ef',
      'database',
      'update',
      '--project',
      infrastructureProject,
      '--startup-project',
      infrastructureProject,
      '--configuration',
      'Release',
    ],
    {
      env: commandEnvironment({
        CROWNPILOT__DATABASE__CONNECTIONSTRING: databaseConnectionString(),
      }),
    },
  );
  console.log('EF Core migrations applied.');
}

function security() {
  assertLocalDatabaseTarget();
  assertDatabaseReachable();
  const migrationState = queryPsql(
    `SELECT CASE WHEN EXISTS (SELECT 1 FROM public."__EFMigrationsHistory" WHERE "MigrationId" = '${preparedSchemaMigration}') THEN 'ready' ELSE 'missing' END;`,
  );
  if (migrationState !== 'ready') {
    failClosed(
      `security SQL requires EF Core migration ${preparedSchemaMigration}; run db:migrate.`,
    );
  }

  const files = sqlFiles(securityDirectory);
  if (files.length === 0) {
    throw new CommandError('No versioned security SQL files were found.');
  }

  runPsql(sqlContents(files));
  console.log(`Security SQL applied (${files.length} file${files.length === 1 ? '' : 's'}).`);
}

function fixtures() {
  assertLocalDatabaseTarget();
  assertDatabaseReachable();
  const securityState = queryPsql(
    `SELECT CASE WHEN EXISTS (SELECT 1 FROM public."__EFMigrationsHistory" WHERE "MigrationId" = '${preparedSchemaMigration}') AND EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'crownpilot') AND EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'crownpilot_runtime') THEN 'ready' ELSE 'missing' END;`,
  );
  if (securityState !== 'ready') {
    failClosed('fixtures require db:security to complete first; run db:security.');
  }

  const files = sqlFiles(fixturesDirectory);
  if (files.length === 0) {
    failClosed(
      'no application fixtures exist yet; 002-06 owns domain tables and sanitized fixtures. No empty success is reported.',
    );
  }

  runPsql(sqlContents(files));
  console.log(`Sanitized fixtures applied (${files.length} file${files.length === 1 ? '' : 's'}).`);
}

async function reset() {
  assertLocalDatabaseTarget();
  assertDockerDaemon();
  compose(['down', '--volumes', '--remove-orphans']);
  await start();
  migrate();
  security();
  fixtures();
  console.log('Local database reset completed: EF migrations -> security SQL -> fixtures.');
}

async function main() {
  if (!['start', 'stop', 'migrate', 'security', 'fixtures', 'reset'].includes(operation)) {
    throw new CommandError(
      `unknown database operation '${operation}'. Use start, stop, migrate, security, fixtures, or reset.`,
      2,
    );
  }

  assertLocalEnvironment();
  assertLocalDatabaseTarget();

  if (operation === 'start') {
    await start();
  } else if (operation === 'stop') {
    stop();
  } else if (operation === 'migrate') {
    migrate();
  } else if (operation === 'security') {
    security();
  } else if (operation === 'fixtures') {
    fixtures();
  } else {
    await reset();
  }
}

main().catch((error) => {
  console.error(error instanceof Error ? error.message : String(error));
  process.exitCode = error?.exitCode ?? 1;
});
