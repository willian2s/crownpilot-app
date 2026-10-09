import { readFile } from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const repositoryRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');

async function text(relativePath) {
  return readFile(path.join(repositoryRoot, relativePath), 'utf8');
}

const domainProject = await text('src/Domain/Domain.csproj');
const applicationProject = await text('src/Application/Application.csproj');
const infrastructureProject = await text('src/Infrastructure/Infrastructure.csproj');
const apiProject = await text('src/Api/Api.csproj');
const domainSource = await text('src/Domain/DomainAssemblyMarker.cs');
const applicationSource = await text('src/Application/ApplicationAssemblyMarker.cs');

const forbiddenInCore = /Firebase|Supabase|Npgsql|EntityFrameworkCore|Microsoft\.AspNetCore|System\.Net\.Http/i;
if (forbiddenInCore.test(domainProject) || forbiddenInCore.test(domainSource)) {
  throw new Error('Domain boundary imports a forbidden framework/provider.');
}
if (forbiddenInCore.test(applicationProject) || forbiddenInCore.test(applicationSource)) {
  throw new Error('Application boundary imports a forbidden framework/provider.');
}

for (const reference of ['..\\Application\\Application.csproj', '..\\Domain\\Domain.csproj']) {
  if (!infrastructureProject.includes(reference)) {
    throw new Error(`Infrastructure is missing ${reference}.`);
  }
}
for (const reference of ['..\\Application\\Application.csproj', '..\\Infrastructure\\Infrastructure.csproj']) {
  if (!apiProject.includes(reference)) {
    throw new Error(`API composition root is missing ${reference}.`);
  }
}
if (!infrastructureProject.includes('Microsoft.EntityFrameworkCore') ||
    !infrastructureProject.includes('Npgsql.EntityFrameworkCore.PostgreSQL')) {
  throw new Error('EF Core/Npgsql packages are not isolated in Infrastructure.');
}

console.log('Layer boundary check passed.');
