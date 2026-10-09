import { withLocalApi } from './process.mjs';

await withLocalApi(async (baseUrl) => {
  const response = await fetch(`${baseUrl}/health/live`);
  if (!response.ok) {
    throw new Error(`Liveness returned HTTP ${response.status}.`);
  }

  const body = await response.text();
  if (body !== 'Healthy') {
    throw new Error(`Liveness returned unexpected body: ${body}`);
  }
});

console.log('Local API smoke passed: /health/live');
