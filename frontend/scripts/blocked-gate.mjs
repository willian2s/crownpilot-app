const gate = process.argv[2] ?? 'unknown';

console.error(`${gate} gate blocked: this bootstrap has no remote database or real authentication flow.`);
console.error('Run this gate after its owning SDD task provisions a disposable fixture environment.');
process.exit(2);
