import { readFile } from 'node:fs/promises';
import path from 'node:path';
import { repositoryRoot, withLocalApi } from './process.mjs';

// Fonte canônica do contrato (spec-first). A API precisa servi-la sem alteração.
const sourcePath = path.join(repositoryRoot, 'api', 'openapi', 'v1.json');
const source = await readFile(sourcePath, 'utf8');

await withLocalApi(async (baseUrl) => {
  const response = await fetch(`${baseUrl}/openapi/v1.json`);
  if (!response.ok) {
    throw new Error(`OpenAPI endpoint returned HTTP ${response.status}.`);
  }

  const served = await response.text();
  if (served !== source) {
    throw new Error('Served /openapi/v1.json differs from api/openapi/v1.json.');
  }

  const document = JSON.parse(served);
  if (typeof document.openapi !== 'string' || !document.openapi.startsWith('3.')) {
    throw new Error('Document does not declare an OpenAPI 3 version.');
  }

  if (!document.paths?.['/api/v1/bootstrap']?.get) {
    throw new Error('Document is missing GET /api/v1/bootstrap.');
  }

  // ADR 005, decisão 6: falha de comunicação com o Firebase na revogação
  // sensível usa authentication_unavailable, nunca provider_unavailable.
  const deleteMe = document.paths?.['/api/v1/me']?.delete;
  const example = deleteMe?.responses?.['503']?.content?.['application/problem+json']?.examples
    ?.authenticationUnavailable?.value;
  if (example?.code !== 'authentication_unavailable') {
    throw new Error('DELETE /api/v1/me does not document 503 authentication_unavailable.');
  }

  const docs = await fetch(`${baseUrl}/docs`);
  if (!docs.ok || !(await docs.text()).includes('/openapi/v1.json')) {
    throw new Error(`Local /docs did not serve Swagger UI for /openapi/v1.json (HTTP ${docs.status}).`);
  }
});

console.log('OpenAPI contract passed: /openapi/v1.json serves api/openapi/v1.json unchanged');
