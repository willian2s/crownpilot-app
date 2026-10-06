const operation = process.argv[2] ?? 'unknown';

console.error(`db:${operation} blocked: PostgreSQL/Supabase setup belongs to task 002-03.`);
console.error('No remote database, schema, migration, fixture, or secret is used by bootstrap.');
process.exit(2);
