import { withLocalApi } from './process.mjs';

await withLocalApi(async (baseUrl) => {
  const response = await fetch(`${baseUrl}/openapi/v1.json`);
  if (!response.ok) {
    throw new Error(`OpenAPI endpoint returned HTTP ${response.status}.`);
  }

  const document = await response.json();
  if (typeof document.openapi !== 'string' || !document.openapi.startsWith('3.')) {
    throw new Error('Generated document does not declare an OpenAPI 3 version.');
  }

  if (!document.paths?.['/api/v1/bootstrap']?.get) {
    throw new Error('Generated document is missing GET /api/v1/bootstrap.');
  }
});

console.log('OpenAPI contract passed: generated /openapi/v1.json');
